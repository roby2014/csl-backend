package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robyzzz/csl-backend/config"
	"github.com/robyzzz/csl-backend/model"
	"github.com/robyzzz/csl-backend/utils"
	"github.com/solovev/steam_go"
)

// GET /api/player/{steamid}
// Returns {steamid}'s steam data or 404 if not found
func GetSteamUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	player, err := model.GetSteamUser(mux.Vars(r)["steamid"])
	if err != nil {
		utils.APIErrorRespond(w, utils.ErrorResponse{Code: http.StatusNotFound, ErrorMsg: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(player)
}

// Not an API function
// Called when user logs in with steam
func CreateSteamUser(user *steam_go.PlayerSummaries) error {
	return model.CreateSteamUser(utils.PlayerSummariesToSteamUser(user))
}

// Not an API function
// Called when we want to update our logged user steam data using his session ID
// Updates steam data or internal server error if smt bad happened
func UpdateSteamUser(w http.ResponseWriter, r *http.Request) {
	steamID := config.GetSessionID(r)

	exists, err := model.DoesSteamUserExist(steamID)
	if err != nil {
		utils.APIErrorRespond(w, utils.ErrorResponse{Code: http.StatusInternalServerError, ErrorMsg: err.Error()})
		return
	}

	updatedUser, err := steam_go.GetPlayerSummaries(steamID, config.STEAM_API_KEY)
	if err != nil {
		utils.APIErrorRespond(w, utils.ErrorResponse{Code: http.StatusInternalServerError, ErrorMsg: err.Error()})
		return
	}

	steamUser := utils.PlayerSummariesToSteamUser(updatedUser)
	var result error
	if exists {
		result = model.UpdateSteamUser(steamUser)
	} else {
		result = model.CreateSteamUser(steamUser)
	}

	if result != nil {
		utils.APIErrorRespond(w, utils.ErrorResponse{Code: http.StatusInternalServerError, ErrorMsg: err.Error()})
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
