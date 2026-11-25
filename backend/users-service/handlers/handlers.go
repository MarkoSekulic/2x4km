package handlers

import (
	"2x4km/users/config"
	"2x4km/users/repo"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"math/rand"
	"net/http"
	"proto/users-service"
	"regexp"
	"shared/proto/files"
	"shared/smtp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type HandlerUtils struct {
	SmtpClient  *smtp.Client
	FilesClient files.FilesServiceClient
}

type UsersHandler struct {
	users.UnimplementedUsersServiceServer
	repo  *repo.UsersRepository
	cfg   *config.UserConfig
	utils HandlerUtils
}

const USER_COLLECTION = "users"
const VERIFICATION_COLLECTION = "verifies"

func InitHandler(repo *repo.UsersRepository, cfg *config.UserConfig, utils HandlerUtils) UsersHandler {
	return UsersHandler{repo: repo, cfg: cfg, utils: utils}
}

func (h UsersHandler) RegisterPreVerify(ctx context.Context, request *users.RegisterRequest) (*users.Empty, error) {

	missing := request.Email == "" || request.FirstName == "" || request.LastName == "" || request.Password == ""
	if missing {
		return nil, status.Error(codes.InvalidArgument, "Missing Fields")
	}

	matched, err := regexp.MatchString(`^[^\s@]+@[^\s@]+\.[^\s@]+$`, request.Email)
	if err != nil || !matched {
		return nil, status.Error(codes.InvalidArgument, "Invalid Email")
	}

	matched, err = regexp.MatchString(`(?s)^(.{0,7}|.{65,}|.*\s.*|[^a-z]*|[^A-Z]*|\D*|[A-Za-z0-9\s]*)$`, request.Password)
	matched = !matched
	if err != nil || !matched {
		return nil, status.Error(codes.InvalidArgument, "Weak Password")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10) //12 is standard
	if err != nil {
		return nil, status.Error(codes.Internal, "Password")
	}

	userDoc := map[string]any{
		"email":        request.Email,
		"firstName":    request.FirstName,
		"lastName":     request.LastName,
		"passwordHash": string(passwordHash),
		"dataLimitMB":  h.cfg.DefaultDataLimitMB,
		"createdAt":    time.Now().UTC(),
	}

	var result bson.M
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = h.repo.Collection(USER_COLLECTION).FindOne(dbCtx, bson.M{"email": request.Email}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&result)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "duplicate email")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, status.Error(codes.Unavailable, "users-db")
	}

	dbCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	res, err := h.repo.Collection(VERIFICATION_COLLECTION).InsertOne(dbCtx, userDoc)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "users-db")
	}
	verificationID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "userID unprocesable")
	}

	err = h.utils.SmtpClient.Send(smtp.Message{
		From:    "noreply@gmail.com",
		To:      []string{request.Email},
		Subject: "2x4km Account Verification",
		Body:    "<a href='localhost:3000/user/verify/" + verificationID.Hex() + "' > Verify account here </a> ",
		IsHTML:  true,
	})
	if err != nil {
		return nil, err
	}

	return &users.Empty{}, nil
}

func (h UsersHandler) Register(ctx context.Context, request *users.VerifyRequest) (*users.Empty, error) {

	verifyId, _ := primitive.ObjectIDFromHex(request.VerifyId)

	session, err := h.repo.Client().StartSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "users-db")
	}

	res, err := session.WithTransaction(ctx, func(dbCtx mongo.SessionContext) (interface{}, error) {
		var resDoc bson.M
		// Verification Collection
		err = h.repo.Collection(VERIFICATION_COLLECTION).FindOne(dbCtx, bson.M{"_id": verifyId}).Decode(&resDoc)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, status.Error(codes.NotFound, "invalid verification id")
			}
			return nil, status.Error(codes.Unavailable, "users-db")
		}
		_, err = h.repo.Collection(VERIFICATION_COLLECTION).DeleteMany(dbCtx, bson.M{"email": resDoc["email"]})

		// Users Collection
		res, err := h.repo.Collection(USER_COLLECTION).InsertOne(dbCtx, resDoc)
		if err != nil {
			return nil, status.Error(codes.Unavailable, "users-db")
		}
		userID, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, status.Error(codes.Internal, "userID unprocesable")
		}

		resDoc["_id"] = userID

		return resDoc, nil
	})
	if err != nil {
		return nil, err
	}

	// CQRS send userID to other services that need them
	userDoc := res.(bson.M)

	h.utils.FilesClient.InnerRootCreate(ctx, &files.InnerRootCreateRequest{
		UserId:    userDoc["_id"].(primitive.ObjectID).Hex(),
		FirstName: userDoc["first_name"].(string),
		LastName:  userDoc["last_name"].(string),
	})
	// TODO MOVE TO LOGIN

	return &users.Empty{}, nil
}

func (h UsersHandler) Login(ctx context.Context, request *users.LoginRequest) (*users.TFATokenResponse, error) {

	var result bson.M
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := h.repo.Collection(USER_COLLECTION).FindOne(dbCtx, bson.M{"email": request.Email}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.Unauthenticated, "wrong email/password")
		}
		return nil, status.Error(codes.Unavailable, "users-db")
	}
	err = bcrypt.CompareHashAndPassword([]byte(result["passwordHash"].(string)), []byte(request.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "wrong email/password")
	}

	passkey := strconv.Itoa(int(rand.Float64() * 1000000))
	TFADoc := bson.M{
		"_type":       "TFA",
		"passkey":     passkey,
		"userId":      result["_id"],
		"email":       result["email"],
		"firstName":   result["firstName"],
		"lastName":    result["lastName"],
		"dataLimitMB": result["dataLimitMB"],
	}
	dbCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	res, err := h.repo.Collection(VERIFICATION_COLLECTION).InsertOne(dbCtx, TFADoc)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "users-db")
	}
	TFAtoken, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "userID unprocesable")
	}

	err = h.utils.SmtpClient.Send(smtp.Message{
		From:    "noreply@gmail.com",
		To:      []string{request.Email},
		Subject: "2x4km 2FA Verification",
		Body:    "<h1>" + passkey + "</h1>",
		IsHTML:  true,
	})
	if err != nil {
		return nil, err
	}

	return &users.TFATokenResponse{TFAToken: TFAtoken.Hex()}, nil
}

func (h UsersHandler) TFALogin(ctx context.Context, request *users.TFALoginRequest) (*users.Empty, error) {

	TFAId, _ := primitive.ObjectIDFromHex(request.TFAToken)
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var result bson.M
	err := h.repo.Collection(VERIFICATION_COLLECTION).FindOne(dbCtx, bson.M{"_id": TFAId}).Decode(&result)
	if err != nil {
		return nil, err
	}
	if request.Passkey != result["passkey"].(string) {
		return nil, status.Error(codes.Unauthenticated, "wrong passkey")
	}

	dbCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = h.repo.Collection(VERIFICATION_COLLECTION).DeleteMany(dbCtx, bson.M{"userId": result["userId"], "_type": "TFA"})
	pkHex, err := hex.DecodeString(h.cfg.JWT_PRIVATE_KEY)
	if err != nil {
		return nil, err
	}
	pk, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), pkHex)
	if err != nil {
		return nil, err
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":          result["userId"],
		"email":       result["email"],
		"firstName":   result["firstName"],
		"lastName":    result["lastName"],
		"dataLimitMB": result["dataLimitMB"],
	}).SignedString(pk)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Key")
	}
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	grpc.SendHeader(ctx, metadata.Pairs("set-cookie", cookie.String()))
	return &users.Empty{}, nil
}

func (h UsersHandler) PrePasswordReset(ctx context.Context, request *users.PrePassResetRequest) (*users.Empty, error) {
	var result bson.M
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := h.repo.Collection(USER_COLLECTION).FindOne(dbCtx, bson.M{"email": request.Email}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.Unauthenticated, "")
		}
		return nil, status.Error(codes.Unavailable, "users-db")
	}

	lastPassChangeDate, ok := result["lastPassChangeDate"]
	if !ok || lastPassChangeDate == "" {
		lastPassChangeDate = result["createdAt"]
	}
	if !lastPassChangeDate.(primitive.DateTime).Time().Before(time.Now().UTC().Add(-24 * time.Hour)) {
		return nil, status.Error(codes.FailedPrecondition, "passwords can be changed once per day")
	}

	PassForgotDoc := map[string]any{
		"_type":           "PF",
		"userId":          result["_id"],
		"email":           result["email"],
		"oldPasswordHash": result["passwordHash"],
	}

	dbCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	res, err := h.repo.Collection(VERIFICATION_COLLECTION).InsertOne(dbCtx, PassForgotDoc)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "users-db")
	}
	PassForgotToken, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "userID unprocesable")
	}
	err = h.utils.SmtpClient.Send(smtp.Message{
		From:    "noreply@gmail.com",
		To:      []string{request.Email},
		Subject: "2x4km Account Password Forgot ",
		Body:    "<a href='localhost:3000/user/passForgot/" + PassForgotToken.Hex() + "' > Verify account here </a> ",
		IsHTML:  true,
	})
	if err != nil {
		return nil, err
	}

	return &users.Empty{}, nil
}

func (h UsersHandler) PasswordReset(ctx context.Context, request *users.PassResetRequest) (*users.Empty, error) {

	passForgotId, _ := primitive.ObjectIDFromHex(request.PassForgotToken)

	var result bson.M
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := h.repo.Collection(VERIFICATION_COLLECTION).FindOne(dbCtx, bson.M{"_id": passForgotId}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.Unauthenticated, "")
		}
		return nil, status.Error(codes.Unavailable, "users-db")
	}

	match := request.Email == result["email"]
	if !match {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	err = bcrypt.CompareHashAndPassword([]byte(result["oldPasswordHash"].(string)), []byte(request.NewPassword))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "password can't be the same as previous")
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 10) //12 is standard
	if err != nil {
		return nil, status.Error(codes.Internal, "Password")
	}

	session, err := h.repo.Client().StartSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "users-db")
	}

	_, err = session.WithTransaction(ctx, func(dbCtx mongo.SessionContext) (interface{}, error) {
		err = h.repo.Collection(USER_COLLECTION).FindOneAndUpdate(dbCtx, &bson.M{"_id": result["userId"]}, &bson.M{"passwordHash": newPasswordHash, "lastPassChangeDate": time.Now().UTC()}).Err()
		if err != nil {
			return nil, status.Error(codes.Unavailable, "users-db")
		}
		_, err = h.repo.Collection(VERIFICATION_COLLECTION).DeleteMany(dbCtx, bson.M{"userId": result["userId"], "_type": "PF"})
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return &users.Empty{}, err
}
