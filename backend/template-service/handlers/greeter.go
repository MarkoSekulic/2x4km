package handlers

import (
	"context"
	"fmt"
	"proto/template-service"
)

type TemplateHandler struct {
	template.UnimplementedTemplateServiceServer
}

func (h TemplateHandler) Greet(ctx context.Context, request *template.TemplateRequest) (*template.TemplateResponse, error) {
	return &template.TemplateResponse{
		Greeting: fmt.Sprintf("Hi %s!", request.Name),
	}, nil
}
