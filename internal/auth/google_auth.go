package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/harmelson/tocouaboa-portfolio/internal/config"
	"google.golang.org/api/idtoken"
)

type GoogleInfo struct {
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	GoogleID      string `json:"google_id"` // Campo "sub" do token
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Exp           int64  `json:"exp"`
}

func ValidateGoogleToken(ctx context.Context, token string) (GoogleInfo, error) {
	if token == "" {
		return GoogleInfo{}, errors.New("auth token not provided")
	}

	token = strings.TrimPrefix(token, "Bearer ")

	payload, err := idtoken.Validate(ctx, token, config.GetEnv("GOOGLE_CLIENT_ID", ""))
	if err != nil {
		return GoogleInfo{}, errors.New("invalid or expired token")
	}

	googleInfo := googlePayloadToGoogleInfo(payload)

	if !googleInfo.EmailVerified {
		return GoogleInfo{}, errors.New("email not verified")
	}

	return googleInfo, nil
}

func googlePayloadToGoogleInfo(payload *idtoken.Payload) GoogleInfo {
	email, _ := payload.Claims["email"].(string)

	emailVerified := false
	if ev, ok := payload.Claims["email_verified"].(bool); ok {
		emailVerified = ev
	}

	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	return GoogleInfo{
		Iss:           payload.Issuer,
		Aud:           payload.Audience,
		GoogleID:      payload.Subject,
		Email:         email,
		EmailVerified: emailVerified,
		Name:          name,
		Picture:       picture,
		Exp:           payload.Expires,
	}
}
