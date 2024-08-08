package usecases

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

type SubscribeUseCase struct {
	q           *pgstore.Queries
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	upgrader    websocket.Upgrader
	mu          *sync.Mutex
}

func NewSubscribeUseCase(q *pgstore.Queries, upgrader websocket.Upgrader, mu *sync.Mutex,
	subscribers map[string]map[*websocket.Conn]context.CancelFunc) *SubscribeUseCase {
	return &SubscribeUseCase{
		q:           q,
		subscribers: subscribers,
		upgrader:    upgrader,
		mu:          mu,
	}
}

func (u *SubscribeUseCase) Execute(w http.ResponseWriter, r *http.Request) {
	rawRoomID := chi.URLParam(r, "room_id")
	roomID, err := uuid.Parse(rawRoomID)

	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	_, err = NewGetRoomByIdUseCase(u.q, r.Context()).Execute(roomID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	c, err := u.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("failed to upgrade connection", "error", err)
		http.Error(w, "failed to upgrade to ws connection", http.StatusBadRequest)
	}

	defer c.Close()

	ctx, cancel := context.WithCancel((r.Context()))

	u.mu.Lock()

	if _, ok := u.subscribers[rawRoomID]; !ok {
		u.subscribers[rawRoomID] = make(map[*websocket.Conn]context.CancelFunc)
	}

	slog.Info("new client connected", "room_id", rawRoomID, "cliend_ip", r.RemoteAddr)
	u.subscribers[rawRoomID][c] = cancel
	u.mu.Unlock()

	<-ctx.Done()

	u.mu.Lock()
	delete(u.subscribers[rawRoomID], c)
	u.mu.Unlock()
}
