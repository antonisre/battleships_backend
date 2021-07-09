package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

// User Struct
type Autopilot struct {
	gorm.Model
	GameId   uint `gorm:"size:100;not null;UNIQUE_INDEX:unq_autopilot;"`
	PlayerId uint `gorm:"size:100; not null;UNIQUE_INDEX:unq_autopilot;"`
}

// Register a new user
func (autopilot *Autopilot) InsertAutopilotData(db *gorm.DB) (*Autopilot, error, uint) {
	if err := db.Debug().Create(&autopilot).Error; err != nil {

		if strings.Contains(err.Error(), "autopilots_player_id_users_id_foreign") {
			return nil, errors.New("error.unknown-user-id"), autopilot.PlayerId
		} else if strings.Contains(err.Error(), "autopilots_game_id_games_id_foreign") {
			return nil, errors.New("error.unknown-game-id"), autopilot.GameId
		}
	}
	return autopilot, nil, 0
}

// Register a new user
func (autopilot *Autopilot) CheckAutopilotData(playerId uint, gameId string, db *gorm.DB) (*Autopilot, error) {
	if err := db.Table("autopilots").
		Where("player_id = ? AND game_id = ?", playerId, gameId).Find(&autopilot).Error; err != nil {
		return nil, err
	}
	return autopilot, nil
}
