package repo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Root struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"userId"`
	RootFolder primitive.ObjectID `bson:"rootFolder"`
}

type Folder struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"`
	Parent   *primitive.ObjectID  `bson:"parent"`
	Children []primitive.ObjectID `bson:"children"`
	Files    []File               `bson:"files"`

	Created time.Time `bson:"created"`
	Edited  time.Time `bson:"edited"`
	Version int       `bson:"_ver"`
}

type File struct {
	Parent  primitive.ObjectID `bson:"parent"`
	OSID    string             `bson:"osid"`

	Created time.Time `bson:"created"`
	Edited  time.Time `bson:"edited"`
	Version int       `bson:"_ver"`
}
