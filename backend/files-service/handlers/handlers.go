package handlers

import (
	"2x4km/files/repo"
	"context"
	"fmt"
	"proto/files-service"
)

type FilesHandler struct {
	files.UnimplementedFilesServiceServer
	repo *repo.FilesRepository
}

func InitHandler(repo *repo.FilesRepository) FilesHandler {
	return FilesHandler{repo: repo}
}

func (h FilesHandler) Greet(ctx context.Context, request *files.FilesRequest) (*files.FilesResponse, error) {
	return &files.FilesResponse{
		Greeting: fmt.Sprintf("Hi from Files %s!", request.Name),
	}, nil
}
