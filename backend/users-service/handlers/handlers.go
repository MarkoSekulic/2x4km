package handlers

import (
	"2x4km/users/repo"
	"context"
	"fmt"
	"proto/users-service"
)

type UsersHandler struct {
	users.UnimplementedUsersServiceServer
	repo *repo.UsersRepository
}

func InitHandler(repo *repo.UsersRepository) UsersHandler {
	return UsersHandler{repo: repo}
}

func (h UsersHandler) Greet(ctx context.Context, request *users.UsersRequest) (*users.UsersResponse, error) {
	return &users.UsersResponse{
		Greeting: fmt.Sprintf("Hi from Users %s!", request.Name),
	}, nil
}
