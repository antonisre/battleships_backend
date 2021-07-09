package models

import (
	"strconv"

	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

// User Struct
type Battleship struct {
	Lives    int    `gorm:"size:100;not null;"`
	Point1   string `gorm:"size:100;not null;d"`
	Point2   string `gorm:"size:100;"`
	Point3   string `gorm:"size:100;"`
	Point4   string `gorm:"size:100;"`
	GameId   uint   `gorm:"size:100;not null"`
	PlayerId uint   `gorm:"size:100; not null"`
	Id       uint   `gorm:"size:100; not null;autoIncrement;primaryKey;"`
}

// Create new game
func (battleShip *Battleship) CreateShips(data []Battleship, db *gorm.DB) (*Battleship, error) {
	conv := make([]interface{}, len(data))

	for i, value := range data {
		conv[i] = value
	}

	err := gormbulk.BulkInsert(db, conv, 20)

	if err != nil {
		return nil, err
	}

	return battleShip, nil
}

func (battleShip *Battleship) GenerateShipData(battleShips [][]string, playerId uint, gameId uint, db *gorm.DB) (*Battleship, error) {
	var battleshipData []Battleship

	for _, ship := range battleShips {
		var point2, point3, point4 string
		lives, _ := strconv.Atoi(ship[0])

		switch len(ship) {
		case 5:
			point4 = ship[4]
			point3 = ship[3]
			point2 = ship[2]
			break
		case 4:
			point4 = ""
			point3 = ship[3]
			point2 = ship[2]
			break
		case 3:
			point4 = ""
			point3 = ""
			point2 = ship[2]
			break
		case 2:
			point4 = ""
			point3 = ""
			point2 = ""
			break
		}

		battleshipData = append(battleshipData, Battleship{
			Lives:    lives,
			Point1:   ship[1],
			Point2:   point2,
			Point3:   point3,
			Point4:   point4,
			PlayerId: playerId,
			GameId:   gameId,
		})
	}

	_, err := battleShip.CreateShips(battleshipData, db)
	if err != nil {
		return nil, err
	}

	return battleShip, nil
}

func (battleShip *Battleship) UpdateBattleShipLives(id uint, db *gorm.DB) (*Battleship, error) {
	err := db.Table("battleships").Where("id = ?", id).Update("lives", gorm.Expr("lives- ?", 1)).Error
	if err != nil {
		return nil, err
	}
	return battleShip, nil
}

func (battleShip *Battleship) AvailableShips(playerId, gameId string, ships *[]Battleship, db *gorm.DB) (*[]Battleship, error) {
	err := db.Table("battleships").Where("player_id = ? AND game_id = ? AND lives > 0", playerId, gameId).
		Find(&ships).Error

	if err != nil {
		return nil, err
	}

	return ships, nil
}
