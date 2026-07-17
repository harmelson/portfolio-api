package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/handlers"
	httpMiddleware "github.com/harmelson/tocouaboa-portfolio/internal/http/middleware"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/routes"
	"github.com/harmelson/tocouaboa-portfolio/internal/service"
)

type Handlers struct {
	User         *handlers.UserHandler
	Subscription *handlers.SubscriptionHandler
}

type Dependencies struct {
	Handlers    *Handlers
	UserService service.UserService
}

func NewRouter(deps *Dependencies) *chi.Mux {
	h := deps.Handlers
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Get("/health", healthCheck)
	r.Get("/ready", readinessCheck)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(httpMiddleware.BodyLimit)

		routes.UserRoutes(r, h.User)
		routes.SubscriptionRoutes(r, h.Subscription, deps.UserService)
	})

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readinessCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
