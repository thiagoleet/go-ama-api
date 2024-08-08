package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/thiagoleet/go-ama-api/internal/api/entity"
	"github.com/thiagoleet/go-ama-api/internal/api/usecases"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

// TODO: change to interface
type apiHandler struct {
	q           *pgstore.Queries
	r           *chi.Mux
	upgrader    websocket.Upgrader
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	mu          *sync.Mutex
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
		q: q,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Only in dev environment!!
				return true
			},
		},
		subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		mu:          &sync.Mutex{},
	}

	r := chi.NewRouter()

	// Adding middlewares
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	// Adding CORS
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Adding Web Socket
	r.Get("/subscribe/{room_id}", a.handleSubscribe)

	// Adding routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", a.handleCreateRoom)
			r.Get("/", a.handleGetRooms)

			r.Route("/{room_id}/messages", func(r chi.Router) {
				r.Get("/", a.handleGetRoomMessages)
				r.Post("/", a.handleCreateRoomMessage)
			})

			r.Route("/{message_id}", func(r chi.Router) {
				r.Get("/", a.handleGetRoomMessage)
				r.Patch("/react", a.handleReactToMessage)
				r.Delete("/react", a.handleRemoveReactFromMessage)
				r.Patch("/answer", a.handleMarkMessageAsAnswered)
			})
		})

	})

	a.r = r

	return a
}

func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var body usecases.CreateRoomInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	u := usecases.NewCreateRoomUseCase(h.q, r.Context())

	response, err := u.Execute(body)

	if err != nil {
		slog.Error("failed to insert room", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request) {
	rawRoomID := chi.URLParam(r, "room_id")
	roomID, err := uuid.Parse(rawRoomID)

	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	var body usecases.CreateRoomMessageInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	u := usecases.NewCreateRoomMessageUseCase(h.q, r.Context())

	response, err := u.Execute(body, roomID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)

	notifyClients := usecases.NewNotifyClientsUseCase(h.mu, h.subscribers).Execute

	go notifyClients(entity.Message{
		Kind:   entity.MessageKindMessageCreated,
		RoomId: rawRoomID,
		Value: entity.MessageMessageCreated{
			ID:      response.ID,
			Message: body.Message,
		},
	})

}

func (h apiHandler) handleGetRooms(w http.ResponseWriter, r *http.Request) {
	u := usecases.NewGetRoomsUseCase(h.q, r.Context())

	response, err := u.Execute()

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("no rooms found", "error", err)
			data, _ := json.Marshal(usecases.GetRoomsResponse{
				Rooms: []entity.RoomDTO{},
				Total: 0,
			})
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(data)

			return
		}

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)

}

func (h apiHandler) handleGetRoomMessages(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleRemoveReactFromMessage(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleMarkMessageAsAnswered(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleSubscribe(w http.ResponseWriter, r *http.Request) {}
