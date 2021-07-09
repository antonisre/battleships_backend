package config

import (
	"battleships_backend/src/models"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// Connect to database
func Connect(DbUser, DbPassword, DbName string) {
	var err error
	DBURI := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbName)

	DB, err = gorm.Open("mysql", DBURI)
	if err != nil {
		fmt.Println("Failed connecting to database")
		panic(err)
	}
	// Migrate the models
	DB.Debug().AutoMigrate(
		&models.User{},
		&models.Game{},
		&models.Battleship{},
		models.Shot{},
		models.Autopilot{},
	)

	//foreign keys
	defer DB.Model(&models.Game{}).AddForeignKey("player_id", "users(id)", "RESTRICT", "RESTRICT")
	defer DB.Model(&models.Game{}).AddForeignKey("opponent_id", "users(id)", "RESTRICT", "RESTRICT")
	defer DB.Model(&models.Game{}).AddForeignKey("player_turn", "users(id)", "RESTRICT", "RESTRICT")
	defer DB.Model(&models.Battleship{}).AddForeignKey("player_id", "users(id)", "RESTRICT", "RESTRICT")
	defer DB.Model(&models.Battleship{}).AddForeignKey("game_id", "games(id)", "RESTRICT", "RESTRICT")
	defer DB.Model(&models.Shot{}).AddForeignKey("player_id", "users(id)", "RESTRICT", "RESTRICT")
	defer DB.Model(&models.Shot{}).AddForeignKey("game_id", "games(id)", "RESTRICT", "RESTRICT")
	defer DB.Model(&models.Autopilot{}).AddForeignKey("player_id", "users(id)", "RESTRICT", "RESTRICT")
	defer DB.Model(&models.Autopilot{}).AddForeignKey("game_id", "games(id)", "RESTRICT", "RESTRICT")
}
