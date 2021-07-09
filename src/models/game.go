package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// User Struct
type Game struct {
	gorm.Model
	PlayerTurn uint   `gorm:"size:100;not null;"`
	Status     string `gorm:"size:100;not null;default:'IN_PROGRESS';UNIQUE_INDEX:unq_game;"`
	PlayerId   uint   `gorm:"size:100;not null;UNIQUE_INDEX:unq_game;" json:"player_id"`
	OpponentId uint   `gorm:"size:100;not null;UNIQUE_INDEX:unq_game;" json:"opponent_id"`
	Won        string `gorm:"size:100"`
}

type GameResponse struct {
	PlayerTurn uint
	PlayerId   uint
	OpponentId uint
	Point1     string
	Point2     string
	Point3     string
	Point4     string
	Id         uint
	Won        string
	Lives      int
}

type UserGames struct {
	Id         uint   `json:"game_id"`
	OpponentId uint   `json:"opponent_id"`
	Status     string `json:"status"`
	PlayerId   uint   `json:"-"`
	Won        string `json:"-"`
}

// Create new game
func (game *Game) CreateGame(db *gorm.DB) (*Game, error) {
	if err := db.Debug().Create(&game).Error; err != nil {
		return nil, err
	}
	return game, nil
}

// Get game data
func (game *Game) GameDetails(gameId, playerId string, gameResponse *[]GameResponse, db *gorm.DB) (*[]GameResponse, error, string) {
	var user User
	err := db.Table("games").
		Select("battleships.lives, battleships.point1, battleships.point2, battleships.point3, battleships.point4, battleships.id, games.*, users.id AS userId").
		Joins("left join battleships on games.id = battleships.game_id left join users on games.player_id = users.id").
		Where("battleships.player_id = ? AND battleships.game_id = ? AND games.status = 'IN_PROGRESS'", playerId, gameId).
		Find(gameResponse).Error

	if err != nil {
		return nil, err, ""
	}

	if len(*gameResponse) == 0 {
		db.Table("users").Where("id = ?", playerId).Find(&user)

		if user.Name == "" {
			return nil, errors.New("error.unknown-user-id"), "player-" + playerId
		} else {
			return nil, errors.New("error.unknown-game-id"), gameId
		}
	}
	return gameResponse, nil, ""
}

// Get all user's games
func (game *Game) GetAllUserGames(playerId string, begin, limit int, userGames *[]UserGames, db *gorm.DB) (*[]UserGames, error, int) {
	var user User
	err := db.Table("games").Where("player_id = ? OR opponent_id = ?", playerId, playerId).Offset(begin).Limit(limit).Find(&userGames).Error

	if len(*userGames) == 0 {
		db.Table("users").Where("id = ?", playerId).Find(&user)

		if user.Name == "" {
			return nil, err, http.StatusNotFound
		} else {
			return nil, err, http.StatusNoContent
		}
	} else if err != nil {
		return nil, err, 0
	}

	return userGames, err, 0
}

func (user *UserGames) CompareOpponentId(userGames *[]UserGames, playerId string) []UserGames {
	slice := make([]UserGames, 0, 0)
	parsedPlayerId, _ := strconv.ParseUint(playerId, 10, 32)

	for _, game := range *userGames {

		if uint64(game.OpponentId) == parsedPlayerId {
			game.OpponentId = game.PlayerId
		}

		//status checkup
		if game.Won != "" && game.Won == playerId {
			game.Status = "WON"
		} else if game.Won != "" {
			game.Status = "LOST"
		}

		slice = append(slice, game)
	}
	return slice
}

func (game *Game) FindGameById(gameId, playerId string, db *gorm.DB) (*Game, error, string, bool) {
	var user User
	var gameCheck Game
	err := db.Table("games").
		Where("(player_id = ? OR opponent_id =  ?) AND id = ? AND games.status = 'IN_PROGRESS'", playerId, playerId, gameId).
		Find(&game).Error

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			db.Table("users").Where("id = ?", playerId).Find(&user)
			db.Table("games").Where("id = ?", gameId).Find(&gameCheck)
			if user.Name == "" {
				return nil, errors.New("error.unknown-user-id"), fmt.Sprintf("player-%s", playerId), false
			} else if gameCheck.Status == "" {
				return nil, errors.New("error.unknown-game-id"), gameId, false
			} else {
				return nil, nil, "", true
			}
		}
		return nil, err, "", false
	}
	return game, nil, "", false
}

func (game *Game) UpdateGameData(gameId string, gameData *Game, db *gorm.DB) (*Game, error) {
	err := db.Table("games").Where("id = ?", gameId).Update(&gameData).Error
	if err != nil {
		return nil, err
	}
	return game, nil
}

// Get game data
func (game *Game) GetUserBattleships(gameId, playerId string, gameResponse *[]GameResponse, db *gorm.DB) (*[]GameResponse, error, string) {
	var user User
	err := db.Table("games").
		Select("battleships.lives, battleships.point1, battleships.point2, battleships.point3, battleships.point4, battleships.id, games.*, users.id AS userId").
		Joins("left join battleships on games.id = battleships.game_id left join users on games.player_id = users.id").
		Where("battleships.player_id = ? AND battleships.game_id = ? AND battleships.lives > 0 AND games.status = 'IN_PROGRESS'", playerId, gameId).
		Find(gameResponse).Error

	if err != nil {
		return nil, err, ""
	}

	if len(*gameResponse) == 0 {
		db.Table("users").Where("id = ?", playerId).Find(&user)

		if user.Name == "" {
			return nil, errors.New("error.unknown-user-id"), "player-" + playerId
		} else {
			return nil, errors.New("error.unknown-game-id"), gameId
		}
	}
	return gameResponse, nil, ""
}
