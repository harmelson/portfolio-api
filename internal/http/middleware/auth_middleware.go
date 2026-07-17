package middleware

import (
	"net/http"

	"github.com/harmelson/tocouaboa-portfolio/internal/auth"
	"github.com/harmelson/tocouaboa-portfolio/internal/contextutil"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/response"
	"github.com/harmelson/tocouaboa-portfolio/internal/service"
)

func AuthMiddleware(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("x-tiger-auth")

			googleInfos, err := auth.ValidateGoogleToken(r.Context(), token)
			if err != nil {
				if !auth.DevAuthEnabled() {
					response.Error(w, http.StatusUnauthorized, "unauthorized")
					return
				}

				devToken := r.Header.Get("x-dev-auth")
				googleInfos, err = auth.ValidateDevToken(devToken)
				if err != nil {
					response.Error(w, http.StatusUnauthorized, "unauthorized")
					return
				}
			}

			user, err := userService.GetByGoogleID(r.Context(), googleInfos.GoogleID)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "user not found")
				return
			}

			ctx := contextutil.WithUserID(r.Context(), user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
