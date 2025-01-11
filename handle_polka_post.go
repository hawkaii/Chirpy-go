package main

import (
	"encoding/json"
	"github.com/hawkaii/Chirpy-go/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerPolkaUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find Api key in header")
		return
	}

	if key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Invalid event")
		return
	}

	err = cfg.DB.UpgradeToChirpyRed(params.Data.UserID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
