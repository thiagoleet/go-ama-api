package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type RemoveReactFromMessageUseCase struct {
	q   *pgstore.Queries
	ctx context.Context
}

type RemoveReactFromMessageUseCaseResponse struct {
	ReactionsCount int64  `json:"reactions_count"`
	MessageID      string `json:"message_id"`
}

func NewRemoveReactFromMessageUseCase(queries *pgstore.Queries, context context.Context) *RemoveReactFromMessageUseCase {
	return &RemoveReactFromMessageUseCase{
		q:   queries,
		ctx: context,
	}
}

func (u *RemoveReactFromMessageUseCase) Execute(messageID uuid.UUID) (*ReactToMessageUseCaseResponse, error) {
	_, err := u.q.GetMessage(u.ctx, messageID)

	if err != nil {
		return nil, err
	}

	reactions_count, err := u.q.RemoveReactionFromMessage(u.ctx, messageID)

	if err != nil {
		return nil, err
	}

	response := ReactToMessageUseCaseResponse{
		ReactionsCount: reactions_count,
		MessageID:      messageID.String(),
	}

	return &response, nil
}
