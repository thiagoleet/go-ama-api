package usecases

import (
	"context"
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/thiagoleet/go-ama-api/internal/api/entity"
)

type NotifyClientsUseCase struct {
	mu          *sync.Mutex
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
}

func NewNotifyClientsUseCase(mu *sync.Mutex,
	subscribers map[string]map[*websocket.Conn]context.CancelFunc) *NotifyClientsUseCase {
	return &NotifyClientsUseCase{
		mu:          mu,
		subscribers: subscribers,
	}
}

func (u *NotifyClientsUseCase) Execute(msg entity.Message) {
	u.mu.Lock()
	defer u.mu.Unlock()

	subscribers, ok := u.subscribers[msg.RoomId]
	if !ok || len(subscribers) == 0 {
		return
	}

	for conn, cancel := range subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			slog.Error("failed to send message to client", "error", err)
			cancel()
		}
	}
}
