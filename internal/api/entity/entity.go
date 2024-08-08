package entity

import "github.com/thiagoleet/go-ama-api/internal/store/pgstore"

const (
	MessageKindMessageCreated = "message_created"
)

type Message struct {
	Kind   string `json:"kind"`
	Value  any    `json:"value"`
	RoomId string `json:"-"`
}

type MessageMessageCreated struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type RoomDTO struct {
	ID    string `json:"id"`
	Theme string `json:"theme"`
}

func MapToRoomsDTO(rooms []pgstore.Room) []RoomDTO {
	var dtoList []RoomDTO
	for _, room := range rooms {
		roomDTO := RoomToDTO(room)
		dtoList = append(dtoList, roomDTO)
	}
	return dtoList
}

func RoomToDTO(room pgstore.Room) RoomDTO {
	roomDTO := RoomDTO{
		ID:    room.ID.String(),
		Theme: room.Theme,
	}

	return roomDTO
}
