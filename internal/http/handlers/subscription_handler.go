package handlers

import (
	"net/http"

	"github.com/harmelson/tocouaboa-portfolio/internal/contextutil"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/response"
	"github.com/harmelson/tocouaboa-portfolio/internal/service"
)

type SubscriptionHandler struct {
	subService service.SubscriptionService
}

func NewSubscriptionHandler(subService service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: subService}
}

// GET /api/v1/subscriptions/me — Retorna subscription do usuário autenticado com dados do plano
func (h *SubscriptionHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := contextutil.GetUserID(ctx)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	sub, err := h.subService.GetByUserIDWithPlan(ctx, userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, sub)
}
