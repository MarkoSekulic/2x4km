package repo

import (
	"2x4km/notifications/config"
	"context"
	"time"

	"github.com/gocql/gocql"
)

type NotificationsRepository struct {
	*gocql.Session
}

type Notification struct {
	UserID     gocql.UUID
	NotifiedAt gocql.UUID
	Type       string
	ByUserID   gocql.UUID
	ByUserName string
	FType      string
	FID        gocql.UUID
	FName      string
	AddedFID   gocql.UUID // zero value when type != "added"
	AddedFName string     // "" when type != "added"
}

func InitDB(ctx context.Context, cfg config.DBConfig) (*NotificationsRepository, func(), error) {
	cluster := gocql.NewCluster(cfg.Address)
	cluster.Keyspace = "notifications_db"
	cluster.Timeout = 10 * time.Second

	// cluster.Authenticator = gocql.PasswordAuthenticator{
	// 	Username: cfg.Username,
	// 	Password: cfg.Password,
	// }

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, nil, err
	}

	return &NotificationsRepository{session}, session.Close, nil
}
