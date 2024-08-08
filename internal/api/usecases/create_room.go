package usecases

import (
	"context"

	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type CreateRoomInput struct {
	Theme string `json:"theme"`
}

type CreateRoomResponse struct {
	ID string `json:"id"`
}

type CreateRoomUseCase struct {
	q   *pgstore.Queries
	ctx context.Context
}

func NewCreateRoomUseCase(queries *pgstore.Queries, context context.Context) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		q:   queries,
		ctx: context,
	}
}

func (u *CreateRoomUseCase) Execute(payload CreateRoomInput) (response *CreateRoomResponse, err error) {

	roomID, err := u.q.InsertRoom(u.ctx, payload.Theme)
	if err != nil {
		return nil, err
	}

	data := CreateRoomResponse{ID: roomID.String()}

	return &data, nil
}
