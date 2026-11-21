package repo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *FilesRepository) CreateRoot(ctx context.Context, userID string) (*Root, error) {
	now := time.Now().UTC()

	// Create the user's root folder
	rootFolder := &Folder{
		ID:       primitive.NewObjectID(),
		Parent:   nil,
		Children: []primitive.ObjectID{},
		Files:    []File{},
		Created:  now,
		Edited:   now,
		Version:  1,
	}

	_, err := r.Folders.InsertOne(ctx, rootFolder)
	if err != nil {
		return nil, fmt.Errorf("create root folder: %w", err)
	}

	// Create the root document referencing the folder
	root := &Root{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		RootFolder: rootFolder.ID,
	}

	_, err = r.Roots.InsertOne(ctx, root)
	if err != nil {
		// Rollback: remove the folder if root creation fails
		_, cleanupErr := r.Folders.DeleteOne(
			ctx,
			bson.M{"_id": rootFolder.ID},
		)

		if cleanupErr != nil {
			return nil, fmt.Errorf(
				"create root: %w; cleanup failed: %v",
				err,
				cleanupErr,
			)
		}

		return nil, fmt.Errorf("create root: %w", err)
	}

	return root, nil
}
