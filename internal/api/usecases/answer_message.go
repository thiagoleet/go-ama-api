package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type AnswerMessageUseCase struct {
	q   *pgstore.Queries
	ctx context.Context
}

type AnswerMessageUseCaseResponse struct {
	MessageID string `json:"message_id"`
}

func NewAnswerMessageUseCase(queries *pgstore.Queries, context context.Context) *AnswerMessageUseCase {
	return &AnswerMessageUseCase{
		q:   queries,
		ctx: context,
	}
}

func (u *AnswerMessageUseCase) Execute(messageID uuid.UUID) (*AnswerMessageUseCaseResponse, error) {
	_, err := u.q.GetMessage(u.ctx, messageID)

	if err != nil {
		return nil, err
	}

	err = u.q.MarkMessageAsAnswered(u.ctx, messageID)

	if err != nil {
		return nil, err
	}

	response := AnswerMessageUseCaseResponse{
		MessageID: messageID.String(),
	}

	return &response, nil
}
