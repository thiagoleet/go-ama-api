package usecases

import (
	"context"

	"github.com/thiagoleet/go-ama-api/internal/api/entity"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type GetRoomsUseCase struct {
	q   *pgstore.Queries
	ctx context.Context
}

type GetRoomsResponse struct {
	Rooms []entity.RoomDTO `json:"rooms"`
	Total int              `json:"total"`
}

func NewGetRoomsUseCase(queries *pgstore.Queries, context context.Context) *GetRoomsUseCase {
	return &GetRoomsUseCase{
		q:   queries,
		ctx: context,
	}
}

func (u *GetRoomsUseCase) Execute() (*GetRoomsResponse, error) {
	rooms, err := u.q.GetRooms(u.ctx)

	if err != nil {
		return nil, err
	}

	dtoList := entity.MapToRoomsDTO(rooms)

	response := GetRoomsResponse{
		Rooms: dtoList,
		Total: len(dtoList),
	}

	return &response, nil

}
