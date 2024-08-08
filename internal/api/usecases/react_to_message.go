package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type ReactToMessageUseCase struct {
	q   *pgstore.Queries
	ctx context.Context
}

type ReactToMessageUseCaseResponse struct {
	ReactionsCount int64  `json:"reactions_count"`
	MessageID      string `json:"message_id"`
}

func NewReactToMessageUseCase(queries *pgstore.Queries, context context.Context) *ReactToMessageUseCase {
	return &ReactToMessageUseCase{
		q:   queries,
		ctx: context,
	}
}

func (u *ReactToMessageUseCase) Execute(messageID uuid.UUID) (*ReactToMessageUseCaseResponse, error) {
	_, err := u.q.GetMessage(u.ctx, messageID)

	if err != nil {
		return nil, err
	}

	reactions_count, err := u.q.ReactToMessage(u.ctx, messageID)

	if err != nil {
		return nil, err
	}

	response := ReactToMessageUseCaseResponse{
		ReactionsCount: reactions_count,
		MessageID:      messageID.String(),
	}

	return &response, nil
}
