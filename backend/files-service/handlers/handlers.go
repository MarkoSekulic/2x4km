package handlers

import (
	"2x4km/files/config"
	"2x4km/files/repo"
	"context"
	"errors"
	"shared/jwtauthenticator"
	"strings"
	"time"

	files "proto/files-service"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ROOT_COLLECTION = "root"
const FOLDER_COLLECTION = "folders"
const FILE_COLLECTION = "files"

const maxNameLen = 255

type FilesHandler struct {
	files.UnimplementedFilesServiceServer
	repo *repo.FilesRepository
	cfg  *config.HandlerConfig
}

func InitFilesHandler(repo *repo.FilesRepository, cfg *config.HandlerConfig) FilesHandler {
	return FilesHandler{repo: repo, cfg: cfg}
}

// ---------------------------------------------------------------------------
// Handleri
// ---------------------------------------------------------------------------

func (h FilesHandler) FolderCreate(ctx context.Context, request *files.FolderCreateRequest) (*files.Empty, error) {

	user, err := jwtauthenticator.RetriveValidUser(ctx, h.cfg.JWTPublicKey)
	if err != nil {
		return nil, err
	}

	if err := h.createNode(ctx, FOLDER_COLLECTION, request.Name, request.ParentId, user["id"].(string)); err != nil {
		return nil, err
	}
	return &files.Empty{}, nil
}

func (h FilesHandler) FolderUpdate(ctx context.Context, request *files.FolderUpdateRequest) (*files.Empty, error) {

	idHex := strings.TrimSpace(request.Id)
	if idHex == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing id")
	}
	folderID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	newName := strings.TrimSpace(request.NewName)
	if newName != "" {
		newName, err = validateName(newName)
		if err != nil {
			return nil, err
		}
	}

	newParentHex := strings.TrimSpace(request.NewParentId)
	if newName == "" && newParentHex == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing parent id")
	}

	var newParentID primitive.ObjectID
	if newParentHex != "" {
		newParentID, err = primitive.ObjectIDFromHex(newParentHex)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid new_parent_id")
		}
	}

	now := time.Now().UTC().Truncate(time.Millisecond)

	session, err := h.repo.Client().StartSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "files-db")
	}
	defer session.EndSession(ctx)

	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// DB errors are returned RAW inside the callback (see createNode)
	_, err = session.WithTransaction(txCtx, func(dbCtx mongo.SessionContext) (interface{}, error) {

		// 1. the folder to update
		var folder struct {
			ID     primitive.ObjectID  `bson:"_id"`
			Name   string              `bson:"name"`
			Parent *primitive.ObjectID `bson:"parent"`
			Root   primitive.ObjectID  `bson:"root"`
			Ver    int64               `bson:"_ver"`
		}
		err := h.repo.Collection(FOLDER_COLLECTION).FindOne(dbCtx, bson.M{"_id": folderID}).Decode(&folder)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, status.Error(codes.NotFound, "folder not found")
			}
			return nil, err
		}
		if folder.Parent == nil {
			return nil, status.Error(codes.InvalidArgument, "root folder can not be updated")
		}
		parentID := *folder.Parent

		// new_parent_id equal to the current parent is not a move, only a rename
		moving := newParentHex != "" && newParentID != parentID

		targetParentID := parentID
		if moving {
			targetParentID = newParentID
		}
		targetName := folder.Name
		if newName != "" {
			targetName = newName
		}

		// 2. moving: new parent must exist, share the root and not be inside the folder itself
		if moving {
			var newParent struct {
				Root primitive.ObjectID `bson:"root"`
			}
			err := h.repo.Collection(FOLDER_COLLECTION).
				FindOne(dbCtx, bson.M{"_id": newParentID}, options.FindOne().SetProjection(bson.M{"root": 1})).
				Decode(&newParent)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return nil, status.Error(codes.NotFound, "new parent folder not found")
				}
				return nil, err
			}
			if newParent.Root != folder.Root {
				return nil, status.Error(codes.InvalidArgument, "new parent has a different root")
			}

			// walk up from the new parent to the root folder
			for current := newParentID; ; {
				if current == folder.ID {
					return nil, status.Error(codes.InvalidArgument, "can not move folder into itself")
				}
				var ancestor struct {
					Parent *primitive.ObjectID `bson:"parent"`
				}
				err := h.repo.Collection(FOLDER_COLLECTION).
					FindOne(dbCtx, bson.M{"_id": current}, options.FindOne().SetProjection(bson.M{"parent": 1})).
					Decode(&ancestor)
				if err != nil {
					return nil, err
				}
				if ancestor.Parent == nil {
					break
				}
				current = *ancestor.Parent
			}
		}

		// 3. write to the target parent first: concurrent writes into the same folder conflict here
		// and WithTransaction retries them, so the name check below always sees the other one
		if moving {
			_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
				bson.M{"_id": newParentID},
				bson.M{
					"$push": bson.M{"child_folders": folder.ID},
					"$set":  bson.M{"edited": now},
				},
			)
		} else {
			_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
				bson.M{"_id": parentID},
				bson.M{"$set": bson.M{"edited": now}},
			)
		}
		if err != nil {
			return nil, err
		}

		// 4. name must be unique in the target parent, among folders and files
		// (the folder itself can not match, because its parent or name is changing)
		if moving || targetName != folder.Name {
			filter := bson.M{"parent": targetParentID, "name": targetName}
			for _, c := range []string{FOLDER_COLLECTION, FILE_COLLECTION} {
				n, err := h.repo.Collection(c).CountDocuments(dbCtx, filter, options.Count().SetLimit(1))
				if err != nil {
					return nil, err
				}
				if n > 0 {
					return nil, status.Error(codes.AlreadyExists, "name already exists in folder")
				}
			}
		}

		// 5. moving: take the folder out of the old parent
		if moving {
			var oldParent struct {
				ChildFolders []primitive.ObjectID `bson:"child_folders"`
			}
			err := h.repo.Collection(FOLDER_COLLECTION).
				FindOne(dbCtx, bson.M{"_id": parentID}, options.FindOne().SetProjection(bson.M{"child_folders": 1})).
				Decode(&oldParent)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return nil, status.Error(codes.Internal, "parent folder missing")
				}
				return nil, err
			}

			remaining := bson.A{}
			for _, id := range oldParent.ChildFolders {
				if id != folder.ID {
					remaining = append(remaining, id)
				}
			}
			_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
				bson.M{"_id": parentID},
				bson.M{"$set": bson.M{"child_folders": remaining, "edited": now}},
			)
			if err != nil {
				return nil, err
			}
		}

		// 6. the folder itself
		set := bson.M{"name": targetName, "edited": now, "created_at": now, "_ver": folder.Ver + 1}
		if moving {
			set["parent"] = newParentID
		}
		_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx, bson.M{"_id": folder.ID}, bson.M{"$set": set})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		if _, isStatus := status.FromError(err); isStatus {
			return nil, err
		}
		return nil, status.Error(codes.Unavailable, "files-db")
	}

	return &files.Empty{}, nil
}

func (h FilesHandler) FolderDelete(ctx context.Context, request *files.FolderDeleteRequest) (*files.Empty, error) {

	idHex := strings.TrimSpace(request.Id)
	if idHex == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing id")
	}
	folderID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	now := time.Now().UTC().Truncate(time.Millisecond)

	session, err := h.repo.Client().StartSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "files-db")
	}
	defer session.EndSession(ctx)

	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = session.WithTransaction(txCtx, func(dbCtx mongo.SessionContext) (interface{}, error) {

		// 1. the folder to delete
		var folder struct {
			ID           primitive.ObjectID   `bson:"_id"`
			Parent       *primitive.ObjectID  `bson:"parent"`
			ChildFolders []primitive.ObjectID `bson:"child_folders"`
			Files        []primitive.ObjectID `bson:"files"`
		}
		err := h.repo.Collection(FOLDER_COLLECTION).FindOne(dbCtx, bson.M{"_id": folderID}).Decode(&folder)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, status.Error(codes.NotFound, "folder not found")
			}
			return nil, err
		}
		if folder.Parent == nil {
			return nil, status.Error(codes.InvalidArgument, "root folder can not be deleted")
		}
		if len(folder.ChildFolders) > 0 || len(folder.Files) > 0 {
			return nil, status.Error(codes.FailedPrecondition, "folder not empty")
		}
		parentID := *folder.Parent

		// 2. take the folder out of the parent
		var parent struct {
			ChildFolders []primitive.ObjectID `bson:"child_folders"`
		}
		err = h.repo.Collection(FOLDER_COLLECTION).
			FindOne(dbCtx, bson.M{"_id": parentID}, options.FindOne().SetProjection(bson.M{"child_folders": 1})).
			Decode(&parent)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, status.Error(codes.Internal, "parent folder missing")
			}
			return nil, err
		}

		remaining := bson.A{}
		for _, id := range parent.ChildFolders {
			if id != folder.ID {
				remaining = append(remaining, id)
			}
		}
		_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
			bson.M{"_id": parentID},
			bson.M{"$set": bson.M{"child_folders": remaining, "edited": now}},
		)
		if err != nil {
			return nil, err
		}

		// 3. the folder itself
		_, err = h.repo.Collection(FOLDER_COLLECTION).DeleteOne(dbCtx, bson.M{"_id": folder.ID})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		if _, isStatus := status.FromError(err); isStatus {
			return nil, err
		}
		return nil, status.Error(codes.Unavailable, "files-db")
	}

	return &files.Empty{}, nil
}

func (h FilesHandler) FolderRead(ctx context.Context, request *files.FolderReadRequest) (*files.FolderReadResponse, error) {

	idHex := strings.TrimSpace(request.Id)
	if idHex == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing id")
	}
	folderID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	var folder struct {
		ID           primitive.ObjectID   `bson:"_id"`
		Name         string               `bson:"name"`
		Parent       *primitive.ObjectID  `bson:"parent"`
		Root         primitive.ObjectID   `bson:"root"`
		ChildFolders []primitive.ObjectID `bson:"child_folders"`
		Files        []primitive.ObjectID `bson:"files"`
		CreatedAt    time.Time            `bson:"created_at"`
		Edited       time.Time            `bson:"edited"`
		Ver          int64                `bson:"_ver"`
	}
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = h.repo.Collection(FOLDER_COLLECTION).FindOne(dbCtx, bson.M{"_id": folderID}).Decode(&folder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.NotFound, "folder not found")
		}
		return nil, status.Error(codes.Unavailable, "files-db")
	}

	// the root folder has no parent, so parent_id is empty
	parentID := ""
	if folder.Parent != nil {
		parentID = folder.Parent.Hex()
	}

	childFolders := make([]string, len(folder.ChildFolders))
	for i, id := range folder.ChildFolders {
		childFolders[i] = id.Hex()
	}
	folderFiles := make([]string, len(folder.Files))
	for i, id := range folder.Files {
		folderFiles[i] = id.Hex()
	}

	return &files.FolderReadResponse{
		Id:           folder.ID.Hex(),
		Name:         folder.Name,
		ParentId:     parentID,
		Root:         folder.Root.Hex(),
		ChildFolders: childFolders,
		Files:        folderFiles,
		CreatedAt:    folder.CreatedAt.UTC().Format(time.RFC3339),
		Edited:       folder.Edited.UTC().Format(time.RFC3339),
		Ver:          folder.Ver,
	}, nil
}

func (h FilesHandler) RootFolderRead(ctx context.Context, _ *files.Empty) (*files.FolderReadResponse, error) {

	user, err := jwtauthenticator.RetriveValidUser(ctx, h.cfg.JWTPublicKey)
	if err != nil {
		return nil, err
	}

	userID, ok := user["id"].(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return nil, status.Error(codes.Unauthenticated, "Invalid user identity")
	}

	var rootDoc struct {
		RootFolderID primitive.ObjectID `bson:"root_folder_id"`
	}

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = h.repo.Collection(ROOT_COLLECTION).
		FindOne(
			dbCtx,
			bson.M{"user_id": userID},
			options.FindOne().SetProjection(bson.M{"root_folder_id": 1}),
		).
		Decode(&rootDoc)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.Internal, "root folder missing")
		}

		return nil, status.Error(codes.Unavailable, "files-db")
	}

	folder, err := h.FolderRead(ctx, &files.FolderReadRequest{
		Id: rootDoc.RootFolderID.Hex(),
	})
	if err != nil {
		return nil, err
	}

	folder.Root = ""

	return folder, nil
}

func (h FilesHandler) FileCreate(ctx context.Context, request *files.FileCreateRequest) (*files.Empty, error) {
	user, err := jwtauthenticator.RetriveValidUser(ctx, h.cfg.JWTPublicKey)
	if err != nil {
		return nil, err
	}

	if err := h.createNode(ctx, FILE_COLLECTION, request.Name, request.ParentId, user["id"].(string)); err != nil {
		return nil, err
	}
	return &files.Empty{}, nil
}

func (h FilesHandler) FileUpdate(ctx context.Context, request *files.FileUpdateRequest) (*files.Empty, error) {

	idHex := strings.TrimSpace(request.Id)
	if idHex == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing id")
	}
	fileID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	newName := strings.TrimSpace(request.NewName)
	if newName != "" {
		newName, err = validateName(newName)
		if err != nil {
			return nil, err
		}
	}

	newParentHex := strings.TrimSpace(request.NewParentId)

	var newParentID primitive.ObjectID
	if newParentHex != "" {
		newParentID, err = primitive.ObjectIDFromHex(newParentHex)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid new_parent_id")
		}
	}

	now := time.Now().UTC().Truncate(time.Millisecond)

	session, err := h.repo.Client().StartSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "files-db")
	}
	defer session.EndSession(ctx)

	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// DB errors are returned RAW inside the callback (see createNode)
	_, err = session.WithTransaction(txCtx, func(dbCtx mongo.SessionContext) (interface{}, error) {

		// 1. the file to update
		var file struct {
			ID     primitive.ObjectID `bson:"_id"`
			Name   string             `bson:"name"`
			Parent primitive.ObjectID `bson:"parent"`
			Root   primitive.ObjectID `bson:"root"`
			Ver    int64              `bson:"_ver"`
		}
		err := h.repo.Collection(FILE_COLLECTION).FindOne(dbCtx, bson.M{"_id": fileID}).Decode(&file)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, status.Error(codes.NotFound, "file not found")
			}
			return nil, err
		}
		parentID := file.Parent

		// new_parent_id equal to the current parent is not a move, only a rename
		moving := newParentHex != "" && newParentID != parentID

		targetParentID := parentID
		if moving {
			targetParentID = newParentID
		}
		targetName := file.Name
		if newName != "" {
			targetName = newName
		}

		// 2. moving: new parent must exist and share the root
		if moving {
			var newParent struct {
				Root primitive.ObjectID `bson:"root"`
			}
			err := h.repo.Collection(FOLDER_COLLECTION).
				FindOne(dbCtx, bson.M{"_id": newParentID}, options.FindOne().SetProjection(bson.M{"root": 1})).
				Decode(&newParent)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return nil, status.Error(codes.NotFound, "new parent folder not found")
				}
				return nil, err
			}
			if newParent.Root != file.Root {
				return nil, status.Error(codes.InvalidArgument, "new parent has a different root")
			}
		}

		// 3. write to the target parent first: concurrent writes into the same folder conflict here
		// and WithTransaction retries them, so the name check below always sees the other one
		if moving {
			_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
				bson.M{"_id": newParentID},
				bson.M{
					"$push": bson.M{"files": file.ID},
					"$set":  bson.M{"edited": now},
				},
			)
		} else {
			_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
				bson.M{"_id": parentID},
				bson.M{"$set": bson.M{"edited": now}},
			)
		}
		if err != nil {
			return nil, err
		}

		// 4. name must be unique in the target parent, among folders and files
		// (the file itself can not match, because its parent or name is changing)
		if moving || targetName != file.Name {
			filter := bson.M{"parent": targetParentID, "name": targetName}
			for _, c := range []string{FOLDER_COLLECTION, FILE_COLLECTION} {
				n, err := h.repo.Collection(c).CountDocuments(dbCtx, filter, options.Count().SetLimit(1))
				if err != nil {
					return nil, err
				}
				if n > 0 {
					return nil, status.Error(codes.AlreadyExists, "name already exists in folder")
				}
			}
		}

		// 5. moving: take the file out of the old parent
		if moving {
			var oldParent struct {
				Files []primitive.ObjectID `bson:"files"`
			}
			err := h.repo.Collection(FOLDER_COLLECTION).
				FindOne(dbCtx, bson.M{"_id": parentID}, options.FindOne().SetProjection(bson.M{"files": 1})).
				Decode(&oldParent)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return nil, status.Error(codes.Internal, "parent folder missing")
				}
				return nil, err
			}

			remaining := bson.A{}
			for _, id := range oldParent.Files {
				if id != file.ID {
					remaining = append(remaining, id)
				}
			}
			_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
				bson.M{"_id": parentID},
				bson.M{"$set": bson.M{"files": remaining, "edited": now}},
			)
			if err != nil {
				return nil, err
			}
		}

		// 6. the file itself
		set := bson.M{"name": targetName, "edited": now, "created_at": now, "_ver": file.Ver + 1}
		if moving {
			set["parent"] = newParentID
		}
		_, err = h.repo.Collection(FILE_COLLECTION).UpdateOne(dbCtx, bson.M{"_id": file.ID}, bson.M{"$set": set})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		if _, isStatus := status.FromError(err); isStatus {
			return nil, err
		}
		return nil, status.Error(codes.Unavailable, "files-db")
	}

	return &files.Empty{}, nil
}

func (h FilesHandler) FileDelete(ctx context.Context, request *files.FileDeleteRequest) (*files.Empty, error) {

	idHex := strings.TrimSpace(request.Id)
	if idHex == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing id")
	}
	fileID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	now := time.Now().UTC().Truncate(time.Millisecond)

	session, err := h.repo.Client().StartSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "files-db")
	}
	defer session.EndSession(ctx)

	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = session.WithTransaction(txCtx, func(dbCtx mongo.SessionContext) (interface{}, error) {

		// 1. the file to delete
		var file struct {
			ID     primitive.ObjectID `bson:"_id"`
			Parent primitive.ObjectID `bson:"parent"`
		}
		err := h.repo.Collection(FILE_COLLECTION).FindOne(dbCtx, bson.M{"_id": fileID}).Decode(&file)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, status.Error(codes.NotFound, "file not found")
			}
			return nil, err
		}

		// 2. take the file out of the parent
		var parent struct {
			Files []primitive.ObjectID `bson:"files"`
		}
		err = h.repo.Collection(FOLDER_COLLECTION).
			FindOne(dbCtx, bson.M{"_id": file.Parent}, options.FindOne().SetProjection(bson.M{"files": 1})).
			Decode(&parent)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, status.Error(codes.Internal, "parent folder missing")
			}
			return nil, err
		}

		remaining := bson.A{}
		for _, id := range parent.Files {
			if id != file.ID {
				remaining = append(remaining, id)
			}
		}
		_, err = h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
			bson.M{"_id": file.Parent},
			bson.M{"$set": bson.M{"files": remaining, "edited": now}},
		)
		if err != nil {
			return nil, err
		}

		// 3. the file itself
		_, err = h.repo.Collection(FILE_COLLECTION).DeleteOne(dbCtx, bson.M{"_id": file.ID})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		if _, isStatus := status.FromError(err); isStatus {
			return nil, err
		}
		return nil, status.Error(codes.Unavailable, "files-db")
	}

	return &files.Empty{}, nil
}

func (h FilesHandler) FileRead(ctx context.Context, request *files.FileReadRequest) (*files.FileReadResponse, error) {

	idHex := strings.TrimSpace(request.Id)
	if idHex == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing id")
	}
	fileID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	var file struct {
		ID        primitive.ObjectID `bson:"_id"`
		Name      string             `bson:"name"`
		Parent    primitive.ObjectID `bson:"parent"`
		Root      primitive.ObjectID `bson:"root"`
		OSID      string             `bson:"OSID"`
		CreatedAt time.Time          `bson:"created_at"`
		Edited    time.Time          `bson:"edited"`
		Ver       int64              `bson:"_ver"`
	}
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = h.repo.Collection(FILE_COLLECTION).FindOne(dbCtx, bson.M{"_id": fileID}).Decode(&file)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.NotFound, "file not found")
		}
		return nil, status.Error(codes.Unavailable, "files-db")
	}

	// OSID is empty until the file has an object store entry
	return &files.FileReadResponse{
		Id:        file.ID.Hex(),
		Name:      file.Name,
		ParentId:  file.Parent.Hex(),
		Root:      file.Root.Hex(),
		CreatedAt: file.CreatedAt.UTC().Format(time.RFC3339),
		Edited:    file.Edited.UTC().Format(time.RFC3339),
		Ver:       file.Ver,
		Osid:      file.OSID,
	}, nil
}

func (h FilesHandler) InnerRootCreate(ctx context.Context, request *files.InnerRootCreateRequest) (*files.Empty, error) {

	session, err := h.repo.Client().StartSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "files-db")
	}
	defer session.EndSession(ctx)

	userId, err := primitive.ObjectIDFromHex(request.UserId)
	if err != nil {
		return nil, err
	}

	txCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = session.WithTransaction(txCtx, func(dbCtx mongo.SessionContext) (interface{}, error) {
		res, err := h.repo.Collection(ROOT_COLLECTION).InsertOne(dbCtx, &bson.M{
			"_id":             userId,
			"user_id":         request.UserId,
			"user_first_name": request.FirstName,
			"user_last_name":  request.LastName,
			"root_folder_id":  primitive.NewObjectID(),
		})
		if err != nil {
			return nil, err
		}

		res, err = h.repo.Collection(FOLDER_COLLECTION).InsertOne(dbCtx, &bson.M{
			"name":          "root",
			"parent":        &primitive.Null{},
			"root":          res.InsertedID, // root id
			"child_folders": &bson.A{},
			"files":         &bson.A{},
			"created_at":    time.Now().UTC(),
		})
		if err != nil {
			return nil, err
		}

		_, err = h.repo.Collection(ROOT_COLLECTION).UpdateByID(dbCtx, userId, &bson.M{
			"$set": &bson.M{
				"root_folder_id": res.InsertedID, // root folder id
			},
		})
		if err != nil {
			return nil, err
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &files.Empty{}, nil
}

// ---------------------------------------------------------------------------
// Zajednicka logika
// ---------------------------------------------------------------------------

func validateName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", status.Error(codes.InvalidArgument, "Missing name")
	}
	if len(name) > maxNameLen || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return "", status.Error(codes.InvalidArgument, "Invalid name")
	}
	return name, nil
}

func (h FilesHandler) createNode(ctx context.Context, collection, rawName, parentHex, userID string) error {

	name, err := validateName(rawName)
	if err != nil {
		return err
	}

	parentHex = strings.TrimSpace(parentHex)
	userID = strings.TrimSpace(userID)

	if parentHex == "" && userID == "" {
		return status.Error(codes.InvalidArgument, "Missing parent_id or user_id")
	}

	fromRoot := parentHex == ""
	var parentID primitive.ObjectID
	if !fromRoot {
		parentID, err = primitive.ObjectIDFromHex(parentHex)
		if err != nil {
			return status.Error(codes.InvalidArgument, "Invalid parent_id")
		}
	}

	// Mongo cuva vreme na milisekundu, pa isto zaokruzujemo
	now := time.Now().UTC().Truncate(time.Millisecond)
	nodeID := primitive.NewObjectID()

	// "parent" i "root" se popunjavaju unutar transakcije, kad se zna roditelj.
	doc := bson.M{
		"_id":        nodeID,
		"name":       name,
		"created_at": now,
		"edited":     now,
		"_ver":       1,
	}

	if collection == FILE_COLLECTION {
		doc["OSID"] = ""
	}

	// u koji niz roditelja ide ID deteta
	parentArray := "files"
	if collection == FOLDER_COLLECTION {
		parentArray = "child_folders"
		// prazni nizovi, NE nil: $push nad null poljem puca
		doc["child_folders"] = bson.A{}
		doc["files"] = bson.A{}
	}

	session, err := h.repo.Client().StartSession()
	if err != nil {
		return status.Error(codes.Unavailable, "files-db")
	}
	defer session.EndSession(ctx)

	txCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// VAZNO: unutar callback-a DB greske vracamo SIROVE (bez status.Error).
	// WithTransaction po labeli greske odlucuje da li da ponovi transakciju;
	// status.Error tu labelu brise. Prevodjenje u gRPC gresku radimo ispod.
	_, err = session.WithTransaction(txCtx, func(dbCtx mongo.SessionContext) (interface{}, error) {

		// 1. root folder korisnika (samo kad roditelj nije poslat)
		if fromRoot {
			var rootDoc struct {
				RootFolderID primitive.ObjectID `bson:"root_folder_id"`
			}
			err := h.repo.Collection(ROOT_COLLECTION).
				FindOne(dbCtx, bson.M{"user_id": userID}, options.FindOne().SetProjection(bson.M{"root_folder_id": 1})).
				Decode(&rootDoc)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					// invarijanta: svaki korisnik MORA imati root dokument
					return nil, status.Error(codes.Internal, "root folder missing")
				}
				return nil, err
			}
			parentID = rootDoc.RootFolderID
		}

		// 2. upis u roditelja + citanje njegovog "root" polja
		// (FindOneAndUpdate vraca dokument pre izmene; "root" se ne menja)
		var parentDoc struct {
			Root primitive.ObjectID `bson:"root"`
		}
		err := h.repo.Collection(FOLDER_COLLECTION).FindOneAndUpdate(dbCtx,
			bson.M{"_id": parentID},
			bson.M{
				"$push": bson.M{parentArray: nodeID},
				"$set":  bson.M{"edited": now},
			},
			options.FindOneAndUpdate().SetProjection(bson.M{"root": 1}),
		).Decode(&parentDoc)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				if fromRoot {
					// root_folder_id ukazuje na folder koji ne postoji: povreda invarijante
					return nil, status.Error(codes.Internal, "root folder missing")
				}
				return nil, status.Error(codes.NotFound, "parent folder not found")
			}
			return nil, err
		}

		// 3. jedinstvenost imena u folderu, preko obe kolekcije
		filter := bson.M{"parent": parentID, "name": name}
		for _, c := range []string{FOLDER_COLLECTION, FILE_COLLECTION} {
			n, err := h.repo.Collection(c).CountDocuments(dbCtx, filter, options.Count().SetLimit(1))
			if err != nil {
				return nil, err
			}
			if n > 0 {
				return nil, status.Error(codes.AlreadyExists, "name already exists in folder")
			}
		}

		// 4. samo dete
		doc["parent"] = parentID
		doc["root"] = parentDoc.Root

		if collection == FOLDER_COLLECTION {
			// doc["OSID"] =
		}

		if _, err := h.repo.Collection(collection).InsertOne(dbCtx, doc); err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		// nasa namerna gRPC greska ide klijentu nepromenjena
		if _, isStatus := status.FromError(err); isStatus {
			return err
		}
		// sve ostalo je problem sa bazom; pravu gresku ne saljemo klijentu
		return err
	}

	return nil
}
