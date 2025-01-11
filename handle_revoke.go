package main

import (
	"net/http"

	"github.com/hawkaii/Chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
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

	// Revoke the refresh token
	err = cfg.DB.DeleteRefreshToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
