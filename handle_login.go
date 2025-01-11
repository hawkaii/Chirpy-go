package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hawkaii/Chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token         string `json:"token"`
		Refresh_token string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	defaultExpiration := 3600
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtKey, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}

	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
	}
	refresh_token := hex.EncodeToString(bytes)

	refreshTokenExpiresAt := time.Now().Add(60 * 24 * time.Hour) // 60 days from now
	err = cfg.DB.StoreRefreshToken(user.ID, refresh_token, refreshTokenExpiresAt)

	respondWithJSON(w, http.StatusOK, response{
		User:          User{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed},
		Token:         token,
		Refresh_token: refresh_token,
	})
}
