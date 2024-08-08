package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/thiagoleet/go-ama-api/internal/api/entity"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type GetRoomByIdUseCase struct {
	q   *pgstore.Queries
	ctx context.Context
}

type GetRoomByIdResponse struct {
	Room entity.RoomDTO `json:"room"`
}

func NewGetRoomByIdUseCase(queries *pgstore.Queries, ctx context.Context) *GetRoomByIdUseCase {
	return &GetRoomByIdUseCase{
		q:   queries,
		ctx: ctx,
	}
}

func (u *GetRoomByIdUseCase) Execute(roomID uuid.UUID) (*GetRoomByIdResponse, error) {
	room, err := u.q.GetRoom(u.ctx, roomID)

	if err != nil {
		return nil, err
	}

	response := GetRoomByIdResponse{
		Room: entity.RoomToDTO(room),
	}

	return &response, nil
}
