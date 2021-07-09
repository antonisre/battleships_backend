package models

import (
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

// User Struct
type Shot struct {
	gorm.Model
	Coordinate string `gorm:"size:100;not null;UNIQUE_INDEX:unq_shot;"`
	PlayerId   uint   `gorm:"size:100;not null;UNIQUE_INDEX:unq_shot;"`
	GameId     uint   `gorm:"size:100;not null;UNIQUE_INDEX:unq_shot;"`
	Status     string `gorm:"size:100;not null;"`
}

type UserShots struct {
	PlayerId   uint
	GameId     uint
	Status     string
	Coordinate string
}

func (shot *Shot) FindUserShots(userId string, gameId string, userShots *[]Shot, db *gorm.DB) (*[]Shot, error) {
	err := db.Table("shots").
		Where("player_id = ? AND game_id = ?", userId, gameId).
		Find(userShots).Error

	if err != nil {
		return nil, err
	}

	return userShots, nil
}

func (shot *Shot) SaveShots(data []Shot, db *gorm.DB) (*Shot, error) {
	conv := make([]interface{}, len(data))

	for i, value := range data {
		conv[i] = value
	}

	err := gormbulk.BulkInsert(db, conv, 20)

	if err != nil {
		return nil, err
	}

	return shot, nil
}
