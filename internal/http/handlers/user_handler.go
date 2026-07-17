package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/harmelson/tocouaboa-portfolio/internal/auth"
	"github.com/harmelson/tocouaboa-portfolio/internal/http/response"
	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/harmelson/tocouaboa-portfolio/internal/service"
)

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *UserHandler) GetByGoogleID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get("x-tiger-auth")

	userInfos, err := auth.ValidateGoogleToken(ctx, token)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.service.GetByGoogleID(ctx, userInfos.GoogleID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get("x-tiger-auth")

	userInfos, err := auth.ValidateGoogleToken(ctx, token)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid authentication token")
		return
	}

	user := models.UserDTO{
		Email:     userInfos.Email,
		Name:      userInfos.Name,
		GoogleID:  userInfos.GoogleID,
		PictureID: userInfos.Picture,
		IsActive:  true,
	}

	if err := h.validate.Struct(user); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Create(ctx, &user); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	response.JSON(w, http.StatusCreated, nil)
}
