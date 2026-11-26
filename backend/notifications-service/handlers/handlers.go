package handlers

import (
	"2x4km/notifications/config"
	"2x4km/notifications/repo"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"proto/notifications-service"

	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HandlerUtils struct {
}

type NotificationsHandler struct {
	notifications.UnimplementedNotificationsServiceServer
	repo  *repo.NotificationsRepository
	cfg   *config.NotificationConfig
	utils HandlerUtils
}

func InitHandler(repo *repo.NotificationsRepository, cfg *config.NotificationConfig, utils HandlerUtils) NotificationsHandler {
	return NotificationsHandler{repo: repo, cfg: cfg, utils: utils}
}

func (h NotificationsHandler) Greet(ctx context.Context, request *notifications.Empty) (*notifications.Empty, error) {
	return &notifications.Empty{}, nil
}

func (h NotificationsHandler) GetNotifications(ctx context.Context, request *notifications.Empty) (*notifications.NotificationsResponse, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	md_vals := md.Get("session_token")
	if len(md_vals) == 0 {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	pkHex, err := hex.DecodeString(h.cfg.JWT_PUBLIC_KEY)
	if err != nil {
		return nil, err
	}
	pk, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), pkHex)
	if err != nil {
		return nil, err
	}

	session_token, err := jwt.Parse(
		md_vals[0],
		func(t *jwt.Token) (any, error) { return pk, nil },
		jwt.WithValidMethods([]string{"ES256"}), // P-256 means ES256
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	user := session_token.Claims.(jwt.MapClaims)

	notifsRaw, err := h.repo.Query("SELECT * FROM notifications_by_user WHERE user_id = ? LIMIT 50;", user["id"]).Iter().SliceMap()
	if err != nil {
		return nil, err
	}
	notifs := []*notifications.Notification{}
	for _, notif := range notifsRaw {
		n := &notifications.Notification{
			UserId:     notif["user_id"].(string),
			NotifiedAt: timestamppb.New(notif["notified_at"].(gocql.UUID).Time()),
			Type:       notif["n_type"].(string),
			ByUserId:   notif["by_user_id"].(string),
			ByUserName: notif["by_user_name"].(string),
			FType:      notif["f_type"].(string),
			FId:        notif["f_id"].(string),
			FName:      notif["f_name"].(string),
		}
		if v, ok := notif["added_f_id"].(string); ok {
			n.AddedFId = v
		}
		if v, ok := notif["added_f_name"].(string); ok {
			n.AddedFName = v
		}
		notifs = append(notifs, n)
	}
	return &notifications.NotificationsResponse{Notifications: notifs}, nil
}

func (h NotificationsHandler) Notify(ctx context.Context, request *notifications.Empty) (*notifications.Empty, error) {

	// id := gocql.TimeUUID()
	// b := session.NewBatch(gocql.LoggedBatch)
	// b.Query(`INSERT INTO notifications_by_user (user_id, notified_at, type, by_user_id, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name)
	// 		VALUES (?,?,?,?,?,?,?,?,?,?)`, uid, id, typ, byID, byName, fType, fID, fName, addedID, addedName)
	// // ...same values into the other three tables
	// err := session.ExecuteBatch(b)

	return nil, status.Error(codes.Unimplemented, "method Notify not implemented")
}

// SELECT * FROM notifications_by_user       WHERE user_id = ? LIMIT 50;
// SELECT * FROM notifications_by_user_type  WHERE user_id = ? AND type = 'shared' LIMIT 50;
// SELECT * FROM notifications_by_user_actor WHERE user_id = ? AND by_user_id = ? LIMIT 50;
// SELECT * FROM notifications_by_user_file  WHERE user_id = ? AND f_id = ? LIMIT 50;

// -- pagination: add AND notified_at < ?  (or use PageState)
