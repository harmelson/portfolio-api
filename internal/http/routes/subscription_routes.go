package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/handlers"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/middleware"
	"github.com/harmelson/tocouaboa-portfolio/internal/service"
)

func SubscriptionRoutes(r chi.Router, h *handlers.SubscriptionHandler, userService service.UserService) {
	r.Route("/subscriptions", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(userService))

		r.Get("/me", h.GetMy)
	})
}
