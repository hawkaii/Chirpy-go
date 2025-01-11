package main

import (
	"net/http"
	"time"

	"github.com/hawkaii/Chirpy-go/internal/auth"
)

type User struct {
	ID                 int       `json:"id"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
	IsChirpyRed        bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	// Extract the refresh token from the Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token")
		return
	}

	// Validate the refresh token and retrieve the associated user
	user, err := cfg.DB.ValidateRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Generate a new JWT (access token) for the user
	newToken, err := auth.MakeJWT(user.ID, cfg.jwtKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new token")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	// Respond with the new JWT
	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})

}
