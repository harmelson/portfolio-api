package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/handlers"
)

func UserRoutes(r chi.Router, h *handlers.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/get", h.GetByGoogleID)
		r.Post("/", h.Create)
	})
}
