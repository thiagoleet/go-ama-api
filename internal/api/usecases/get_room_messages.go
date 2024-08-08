package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/thiagoleet/go-ama-api/internal/api/entity"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type GetRoomMessages struct {
	q   *pgstore.Queries
	ctx context.Context
}

type GetRoomMessagesResponse struct {
	Messages []entity.MessageDTO `json:"messages"`
	RoomID   string              `json:"room_id"`
	Total    int64               `json:"total"`
}

func NewGetRoomMessages(queries *pgstore.Queries, ctx context.Context) *GetRoomMessages {
	return &GetRoomMessages{
		q:   queries,
		ctx: ctx,
	}
}

func (u *GetRoomMessages) Execute(roomID uuid.UUID) (*GetRoomMessagesResponse, error) {
	messages, err := u.q.GetRoomMessages(u.ctx, roomID)

	if err != nil {
		return nil, err
	}

	response := GetRoomMessagesResponse{
		Messages: entity.MapToMessagesDTO(messages),
		RoomID:   roomID.String(),
		Total:    int64(len(messages)),
	}

	return &response, nil
}
