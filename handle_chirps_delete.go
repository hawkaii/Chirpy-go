package main

import (
	"github.com/hawkaii/Chirpy-go/internal/auth"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handleChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert subject to int")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if chirp.AuthorID != id {
		respondWithError(w, http.StatusForbidden, "You are not the author of this chirp")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		http.Error(w, "Failed to delete chirp", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
