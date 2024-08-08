package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type CreateRoomMessageUseCase struct {
	q   *pgstore.Queries
	ctx context.Context
}

type CreateRoomMessageResponse struct {
	ID string `json:"id"`
}

type CreateRoomMessageInput struct {
	Message string `json:"message"`
}

func NewCreateRoomMessageUseCase(queries *pgstore.Queries, context context.Context) *CreateRoomMessageUseCase {
	return &CreateRoomMessageUseCase{
		q:   queries,
		ctx: context,
	}
}

func (u *CreateRoomMessageUseCase) Execute(input CreateRoomMessageInput, roomID uuid.UUID) (*CreateRoomMessageResponse, error) {

	_, err := NewGetRoomByIdUseCase(u.q, u.ctx).Execute(roomID)

	if err != nil {
		return nil, err
	}

	messageID, err := u.q.InsertMessage(u.ctx, pgstore.InsertMessageParams{
		RoomID:  roomID,
		Message: input.Message,
	})

	if err != nil {

		return nil, err
	}

	response := CreateRoomMessageResponse{
		ID: messageID.String(),
	}

	return &response, nil

}
