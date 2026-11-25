package jwtauthenticator

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func RetriveValidUser(ctx context.Context, jwtKey string) (jwt.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	md_vals := md.Get("session_token")
	if len(md_vals) == 0 {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	pkHex, err := hex.DecodeString(jwtKey)
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
	return session_token.Claims.(jwt.MapClaims), nil
}
