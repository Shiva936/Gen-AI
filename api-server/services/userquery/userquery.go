package userquery

import (
	"api-server/clients"
	"api-server/dtos"
	"context"
)

type IUserQuery interface {
	AskAnything(ctx *context.Context, req *dtos.Request) (*dtos.Response, error)
}

type UserQuery struct {
	client clients.IOpenAI
}

func NewUserQuery() IUserQuery {
	return &UserQuery{
		client: clients.NewOpenAI(),
	}
}

func (t *UserQuery) AskAnything(ctx *context.Context, req *dtos.Request) (*dtos.Response, error) {

	content, err := t.EmbedQuery(req.Question)
	if err != nil {
		return nil, err
	}

	answer, err := t.client.Query(ctx, req.Question, content, 3)
	if err != nil {
		return nil, err
	}

	return &dtos.Response{
		Answer: answer,
	}, nil
}

func (t *UserQuery) EmbedQuery(string) ([]string, error) {
	
	return nil, nil
}
