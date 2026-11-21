package middleware

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

func ResponseHeaderMiddleware(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if cookies := md.HeaderMD.Get("set-cookie"); len(cookies) > 0 {
		w.Header().Del("Grpc-Metadata-Set-Cookie")
		for _, cookie := range cookies {
			w.Header().Add("Set-Cookie", cookie)
		}
	}
	return nil
}
