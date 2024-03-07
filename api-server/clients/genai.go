package clients

import "context"

type IOpenAI interface {
	Query(ctx *context.Context, question string, content []string, topK int) (string, error)
}

type OpenAI struct {
}

func NewOpenAI() IOpenAI {
	return &OpenAI{}
}

func (t *OpenAI) Query(ctx *context.Context, question string, content []string, topK int) (string, error) {

	return "", nil
}
