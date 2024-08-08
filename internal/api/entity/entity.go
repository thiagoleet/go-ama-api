package entity

import "github.com/thiagoleet/go-ama-api/internal/store/pgstore"

const (
	MessageKindMessageCreated        = "message_created"
	MessageKindMessageReactAdded     = "message_react_added"
	MessageKindMessageReactedRemoved = "message_react_removed"
	MessageKindMessageAnswered       = "message_answered"
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

type MessageMessageReactAdded struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

type MessageMessageReactRemoved struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

type MessageMessageAnswered struct {
	ID string `json:"id"`
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

	if dtoList == nil {
		return []RoomDTO{}
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

type MessageDTO struct {
	ID             string `json:"id"`
	RoomID         string `json:"room_id"`
	Message        string `json:"message"`
	ReactionsCount int64  `json:"reactions_count"`
	Answered       bool   `json:"answered"`
}

func MessageToDTO(message pgstore.Message) MessageDTO {
	return MessageDTO{
		ID:             message.ID.String(),
		RoomID:         message.RoomID.String(),
		Message:        message.Message,
		ReactionsCount: message.ReactionsCount,
		Answered:       message.Answered,
	}
}

func MapToMessagesDTO(messages []pgstore.Message) []MessageDTO {
	var dtoList []MessageDTO
	for _, message := range messages {
		messageDTO := MessageToDTO(message)
		dtoList = append(dtoList, messageDTO)
	}

	if dtoList == nil {
		return []MessageDTO{}
	}

	return dtoList
}
