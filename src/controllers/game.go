package controllers

import (
	"battleships_backend/config"
	"battleships_backend/src/lib"
	"battleships_backend/src/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateGame(response http.ResponseWriter, request *http.Request) {
	game := &models.Game{}
	ship := &models.Battleship{}
	params := mux.Vars(request)
	opponentId, _ := strconv.ParseUint(params["opponent_id"], 10, 32)
	game.OpponentId = uint(opponentId)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	err = json.Unmarshal(body, &game)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	if game.OpponentId == game.PlayerId {
		lib.Error(response, http.StatusBadRequest, "Duplicated id", "")
		return
	}

	//random player turn
	playerIDs := []interface{}{game.OpponentId, game.PlayerId}
	randomPlayerTurn := lib.PickRandomValue(playerIDs).(uint)
	game.PlayerTurn = randomPlayerTurn

	gameData, err := game.CreateGame(config.DB)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			lib.Error(response, http.StatusBadRequest, "Game already in progress", "")
			return
		} else if strings.Contains(err.Error(), "games_player_id") {
			lib.Error(response, http.StatusBadRequest, "error.unknown-user-id", strconv.FormatUint(uint64(game.PlayerId), 10))
			return
		} else {
			lib.Error(response, http.StatusBadRequest, "error.unknown-user-id", strconv.FormatUint(uint64(game.OpponentId), 10))
			return
		}
	}

	//player's board
	battleShips := lib.GeneratePlayerBattleships()
	ship.GenerateShipData(battleShips, game.PlayerId, gameData.ID, config.DB)

	//opponent's board
	battleShipsOpponenOpponent := lib.GeneratePlayerBattleships()
	ship.GenerateShipData(battleShipsOpponenOpponent, game.OpponentId, gameData.ID, config.DB)

	response.Header().Set("Location", fmt.Sprintf("/game/match-%d", gameData.ID))
	lib.Success(response, http.StatusCreated, map[string]string{
		"player_id":   fmt.Sprintf("player-%d", gameData.PlayerId),
		"opponent_id": fmt.Sprintf("player-%d", gameData.OpponentId),
		"game_id":     fmt.Sprintf("match-%d", gameData.ID),
		"starting":    fmt.Sprintf("player-%d", gameData.PlayerTurn),
	})
	return
}

func ViewGameStatus(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	playerId := params["player_id"]
	gameId := params["game_id"]
	game := &models.Game{}
	userShips := &[]models.GameResponse{}
	userShots := &models.Shot{}
	userShotsResponse := &[]models.Shot{}
	var shipMap = make(map[string]bool)
	var shotMap = make(map[string]bool)
	var opponentId uint
	var gameOwner uint
	var playerTurn uint

	type ResponseUser struct {
		PlayerId string `json:"player_id"`
		Board    []string
	}
	type ResponseData struct {
		Self     ResponseUser      `json:"self"`
		Opponent ResponseUser      `json:"opponent"`
		Game     map[string]string `json:"game"`
	}

	_, err, errArg := game.GameDetails(gameId, playerId, userShips, config.DB)

	if err != nil {
		lib.Error(response, http.StatusNotFound, err.Error(), errArg)
		return
	}

	for _, ship := range *userShips {
		shipMap[ship.Point1] = true
		shipMap[ship.Point2] = true
		shipMap[ship.Point3] = true
		shipMap[ship.Point4] = true
		opponentId = ship.OpponentId
		gameOwner = ship.PlayerId
		playerTurn = ship.PlayerTurn
	}

	usersGrid := lib.GenerateStringGrid(shipMap, false)
	shotInfo, err := userShots.FindUserShots(playerId, gameId, userShotsResponse, config.DB)

	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	for _, shot := range *shotInfo {
		if shot.Status == "HIT" || shot.Status == "KILL" {
			shotMap[shot.Coordinate] = true
		} else {
			shotMap[shot.Coordinate] = false
		}
	}

	opponentGrid := lib.GenerateStringGrid(shotMap, true)
	if playerId != strconv.FormatUint(uint64(gameOwner), 10) {
		opponentId = gameOwner
	}

	responseData := &ResponseData{
		Self: ResponseUser{
			PlayerId: playerId,
			Board:    usersGrid,
		},
		Opponent: ResponseUser{
			PlayerId: strconv.FormatUint(uint64(opponentId), 10),
			Board:    opponentGrid,
		},
		Game: map[string]string{
			"player_turn": strconv.FormatUint(uint64(playerTurn), 10),
		},
	}

	lib.Success(response, http.StatusCreated, &responseData)
	return
}

func GetAllGamesByUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	playerId := params["player_id"]
	game := &models.Game{}
	userGame := &models.UserGames{}
	userGames := &[]models.UserGames{}
	page := 1
	limitPerPage := 40
	offset := (page - 1) * limitPerPage

	_, err, responseStatus := game.GetAllUserGames(playerId, offset, limitPerPage, userGames, config.DB)

	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	if responseStatus != 0 {
		lib.Error(response, responseStatus, "", "")
		return
	}

	filteredGames := userGame.CompareOpponentId(userGames, playerId)
	lib.Success(response, http.StatusCreated, filteredGames)
}

func FireShots(response http.ResponseWriter, request *http.Request) {
	type ShotInput struct {
		Salvo []string
	}
	type ResponseData struct {
		Salvo map[string]string `json:"salvo"`
		Game  map[string]uint   `json:"game"`
	}

	var salvo ShotInput
	params := mux.Vars(request)
	playerId := params["player_id"]
	gameId := params["game_id"]

	responseData := ResponseData{
		Salvo: make(map[string]string),
		Game:  make(map[string]uint),
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	err = json.Unmarshal(body, &salvo)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	shotsData, errShots, httpStatus, errorArg := lib.FireShotsLogic(salvo.Salvo, gameId, playerId)

	if errShots != nil {
		lib.Error(response, httpStatus, fmt.Sprint(errShots), errorArg)
		return
	} else if httpStatus >= 400 {
		lib.Error(response, httpStatus, "", "")
		return
	}

	responseData.Salvo = shotsData.Coordinates
	responseData.Game = map[string]uint{
		"player_turn": shotsData.OpponentId,
	}

	if shotsData.Won != "" {
		parsedWon, _ := strconv.Atoi(shotsData.Won)
		responseData.Game = map[string]uint{
			"won": uint(parsedWon),
		}
	}
	lib.Success(response, http.StatusOK, responseData)
}

func TurnAutopilot(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	playerId := params["player_id"]
	gameId := params["game_id"]
	autopilot := &models.Autopilot{}
	game := &models.Game{}
	parsedPlayerId, err := strconv.ParseUint(playerId, 10, 64)
	parsedGameId, err := strconv.ParseUint(gameId, 10, 64)

	autopilot.PlayerId = uint(parsedPlayerId)
	autopilot.GameId = uint(parsedGameId)

	_, err, _, _ = game.FindGameById(gameId, playerId, config.DB)
	if game.Won != "" {
		lib.Error(response, http.StatusMethodNotAllowed, "", "")
		return
	}

	_, err, errArg := autopilot.InsertAutopilotData(config.DB)

	if err != nil {
		lib.Error(response, http.StatusNotFound, err.Error(), fmt.Sprint(errArg))
		return
	}

	if game.PlayerTurn == uint(parsedPlayerId) {
		coordinates := make([]string, 0)
		lib.FireShotsLogic(coordinates, gameId, playerId)
	}

	lib.Success(response, http.StatusNoContent, "")
}
