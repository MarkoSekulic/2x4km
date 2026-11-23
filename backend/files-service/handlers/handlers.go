package handlers

import (
	"2x4km/files/repo"
	"context"
	"proto/files-service"
	"strings"
	"time"

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
}

func InitFilesHandler(repo *repo.FilesRepository) FilesHandler {
	return FilesHandler{repo: repo}
}

// ---------------------------------------------------------------------------
// Handleri
// ---------------------------------------------------------------------------

func (h FilesHandler) FolderCreate(ctx context.Context, request *files.FolderCreateRequest) (*files.Empty, error) {
	if err := h.createNode(ctx, FOLDER_COLLECTION, request.Name, request.ParentId, request.UserId); err != nil {
		return nil, err
	}
	return &files.Empty{}, nil
}

func (h FilesHandler) FileCreate(ctx context.Context, request *files.FileCreateRequest) (*files.Empty, error) {
	if err := h.createNode(ctx, FILE_COLLECTION, request.Name, request.ParentId, request.UserId); err != nil {
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

// createNode pravi fajl ili folder (collection odredjuje sta).
//
// Roditelj se odredjuje ovako:
//   - parentHex nije prazan  -> to je roditelj, userID se IGNORISE (bez provere vlasnistva)
//   - parentHex je prazan    -> roditelj je root folder korisnika userID (iz root kolekcije)
//   - oba prazna             -> InvalidArgument
//
// Sve se desava u jednoj transakciji:
//  1. (ako treba) nalazenje root foldera korisnika
//  2. $push ID-a deteta u roditeljev files[] / folders[] (+ roditeljev "edited")
//  3. provera da ime ne postoji u tom folderu (ni medju fajlovima ni medju folderima)
//  4. insert deteta
//
// Zasto je $push pre provere imena: dve istovremene transakcije koje kreiraju
// u isti folder obe pisu isti roditeljski dokument, pa Mongo jednu prekine sa
// WriteConflict (TransientTransactionError) i WithTransaction je ponovi. Pri
// ponovnom pokusaju provera imena vidi dete one prve, pa duplikat ne moze da prodje.
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

	// roditelj iz zahteva (ako je poslat)
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

	// "parent" u detetu se popunjava tek kad se zna roditelj (unutar transakcije)
	doc := bson.M{
		"_id":     nodeID,
		"name":    name,
		"created": now,
		"edited":  now,
		"_ver":    1,
	}

	// u koji niz roditelja ide ID deteta
	parentArray := "files"
	if collection == FOLDER_COLLECTION {
		parentArray = "folders"
		// prazni nizovi, NE nil: $push nad null poljem puca
		doc["folders"] = bson.A{}
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
				RootFolder primitive.ObjectID `bson:"root_folder"`
			}
			err := h.repo.Collection(ROOT_COLLECTION).
				FindOne(dbCtx, bson.M{"user_id": userID}, options.FindOne().SetProjection(bson.M{"root_folder": 1})).
				Decode(&rootDoc)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// invarijanta: svaki korisnik MORA imati root dokument
					return nil, status.Error(codes.Internal, "root folder missing")
				}
				return nil, err
			}
			parentID = rootDoc.RootFolder
		}

		// 2. upis u roditelja (ujedno proverava da roditelj postoji)
		res, err := h.repo.Collection(FOLDER_COLLECTION).UpdateOne(dbCtx,
			bson.M{"_id": parentID},
			bson.M{
				"$push": bson.M{parentArray: nodeID},
				"$set":  bson.M{"edited": now},
			},
		)
		if err != nil {
			return nil, err
		}
		if res.MatchedCount == 0 {
			if fromRoot {
				// root_folder ukazuje na folder koji ne postoji: povreda invarijante
				return nil, status.Error(codes.Internal, "root folder missing")
			}
			return nil, status.Error(codes.NotFound, "parent folder not found")
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
		return status.Error(codes.Unavailable, "files-db")
	}

	return nil
}
