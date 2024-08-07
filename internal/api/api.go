package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/thiagoleet/go-ama-api/internal/store/pgstore"
)

// TODO: change to interface
type apiHandler struct {
	q *pgstore.Queries
	r *chi.Mux
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
		q: q,
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

func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleGetRooms(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleGetRoomMessages(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleRemoveReactFromMessage(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleMarkMessageAsAnswered(w http.ResponseWriter, r *http.Request) {}

func (h apiHandler) handleSubscribe(w http.ResponseWriter, r *http.Request) {}
