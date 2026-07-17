package auth

import (
	"errors"
	"os"
)

// DevAuthEnabled verifica se a autenticação de desenvolvimento está ativada.
// Útil para recrutadores/testers rodarem a API sem um token Google real.
func DevAuthEnabled() bool {
	return os.Getenv("DEV_AUTH_ENABLED") == "true"
}

// ValidateDevToken valida um token de desenvolvimento e retorna um GoogleInfo
// fixo para o usuário de teste seedado no banco.
func ValidateDevToken(token string) (GoogleInfo, error) {
	if token == "" {
		return GoogleInfo{}, errors.New("dev auth token not provided")
	}

	expected := os.Getenv("DEV_AUTH_TOKEN")
	if expected == "" {
		expected = "dev-token"
	}

	if token != expected {
		return GoogleInfo{}, errors.New("invalid dev auth token")
	}

	return GoogleInfo{
		GoogleID:      "dev-test-user",
		Email:         "dev@example.com",
		EmailVerified: true,
		Name:          "Test User",
		Picture:       "https://example.com/picture.png",
	}, nil
}
