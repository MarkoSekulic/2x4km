package repo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *FilesRepository) CreateFolder(
	ctx context.Context,
	parentID primitive.ObjectID,
) (*Folder, error) {

	// Check whether parent folder exists
	err := r.Folders.FindOne(
		ctx,
		bson.M{"_id": parentID},
	).Err()

	if err != nil {
		return nil, fmt.Errorf("find parent folder: %w", err)
	}

	now := time.Now().UTC()

	// Create new folder
	folder := &Folder{
		ID:       primitive.NewObjectID(),
		Parent:   &parentID,
		Children: []primitive.ObjectID{},
		Files:    []File{},
		Created:  now,
		Edited:   now,
		Version:  1,
	}

	_, err = r.Folders.InsertOne(ctx, folder)
	if err != nil {
		return nil, fmt.Errorf("insert folder: %w", err)
	}

	// Add new folder ID to parent's children array
	result, err := r.Folders.UpdateOne(
		ctx,
		bson.M{"_id": parentID},
		bson.M{
			"$push": bson.M{
				"children": folder.ID,
			},
			"$set": bson.M{
				"edited": now,
			},
			"$inc": bson.M{
				"_ver": 1,
			},
		},
	)

	if err != nil {
		// Rollback: remove the newly created folder
		_, cleanupErr := r.Folders.DeleteOne(
			ctx,
			bson.M{"_id": folder.ID},
		)

		if cleanupErr != nil {
			return nil, fmt.Errorf(
				"update parent: %w; cleanup failed: %v",
				err,
				cleanupErr,
			)
		}

		return nil, fmt.Errorf("update parent: %w", err)
	}

	if result.MatchedCount == 0 {
		_, _ = r.Folders.DeleteOne(ctx, bson.M{"_id": folder.ID})
		return nil, fmt.Errorf("parent folder not found")
	}

	return folder, nil
}


func (r *FilesRepository) AddFile(
    ctx context.Context,
    folderID primitive.ObjectID,
    osid string,
) error {

    now := time.Now().UTC()

    file := File{
        Parent:  folderID,
        OSID:    osid,
        Created: now,
        Edited:  now,
        Version: 1,
    }

    result, err := r.Folders.UpdateOne(
        ctx,
        bson.M{"_id": folderID},
        bson.M{
            "$push": bson.M{"files": file},
            "$set":  bson.M{"edited": now},
            "$inc":  bson.M{"_ver": 1},
        },
    )

    if err != nil {
        return fmt.Errorf("add file: %w", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("folder not found")
    }

    return nil
}

