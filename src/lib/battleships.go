package lib

import (
	"battleships_backend/config"
	"battleships_backend/src/models"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var availableColumns = map[int]string{0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J"}
var availableRows = map[int]int{0: 1, 1: 2, 2: 3, 3: 4, 4: 5, 5: 6, 6: 7, 7: 8, 8: 9, 9: 10}

type ResponseData struct {
	Coordinates map[string]string
	OpponentId  uint
	Won         string
}

func ReserveSurroundingSlots(coordinates []string, occupiedFields map[string]bool) {
	for _, value := range coordinates {
		spplitedCoordinates := strings.Split(value, "x")
		currentPointRow, _ := strconv.Atoi(spplitedCoordinates[0])
		currentPointColumn := []rune(spplitedCoordinates[1])[0]

		coordinate1 := fmt.Sprintf("%dx%s", currentPointRow, string(currentPointColumn+1))
		coordinate2 := fmt.Sprintf("%dx%s", currentPointRow, string(currentPointColumn-1))
		coordinate3 := fmt.Sprintf("%dx%s", currentPointRow+1, string(currentPointColumn))
		coordinate4 := fmt.Sprintf("%dx%s", currentPointRow-1, string(currentPointColumn))

		occupiedFields[coordinate1] = true
		occupiedFields[coordinate2] = true
		occupiedFields[coordinate3] = true
		occupiedFields[coordinate4] = true
		occupiedFields[value] = true
	}
}

func CreateBattleship(shipSize int, occupiedFields map[string]bool) []string {
	var coordinates []string
	var tempCoordinates []string
	numberOfStepsDown := 1
	numberOfStepsLeft := 1

	startPointColumn := GetRandomValueInRange(0, 10)
	startPointRow := GetRandomValueInRange2(0, 10)
	startPoint := fmt.Sprintf("%dx%s", availableRows[startPointRow], availableColumns[startPointColumn])

	for len(coordinates) < shipSize {

		//find starting point
		for _, ok := occupiedFields[startPoint]; ok; {
			startPointColumn = GetRandomValueInRange2(0, 10)
			startPointRow = GetRandomValueInRange(0, 10)
			startPoint = fmt.Sprintf("%dx%s", availableRows[startPointRow], availableColumns[startPointColumn])
			_, ok = occupiedFields[startPoint]
		}
		tempCoordinates = append(tempCoordinates, startPoint)

		//check for free slots by row
		for len(tempCoordinates) < shipSize {
			numberOfStepsUp := len(tempCoordinates)

			//check right side
			currentPointRow, isInRange := availableRows[startPointRow+numberOfStepsUp]
			currentPoint := fmt.Sprintf("%dx%s", currentPointRow, availableColumns[startPointColumn])
			_, occupied := occupiedFields[currentPoint]

			if occupied || !isInRange {

				//check left side of x
				currentPointRow, isInRange := availableRows[startPointRow-numberOfStepsDown]
				currentPoint = fmt.Sprintf("%dx%s", currentPointRow, availableColumns[startPointColumn])
				_, occupied := occupiedFields[currentPoint]

				if occupied || !isInRange {
					//cleanup temp cordinates
					tempCoordinates = tempCoordinates[:1]
					break
				} else {
					numberOfStepsDown++
					tempCoordinates = append(tempCoordinates, currentPoint)
				}

			} else {
				tempCoordinates = append(tempCoordinates, currentPoint)
			}
		}

		//check for free slots by column
		for len(tempCoordinates) < shipSize {
			numberOfStepsRight := len(tempCoordinates)

			//check upper side
			currentPointColumn, isInRange := availableColumns[startPointColumn+numberOfStepsRight]
			currentPoint := fmt.Sprintf("%dx%s", availableRows[startPointRow], currentPointColumn)
			_, occupied := occupiedFields[currentPoint]

			if occupied || !isInRange {

				//check lower side
				currentPointColumn, isInRange := availableColumns[startPointColumn-numberOfStepsLeft]
				currentPoint = fmt.Sprintf("%dx%s", availableRows[startPointRow], currentPointColumn)
				_, occupied := occupiedFields[currentPoint]

				if occupied || !isInRange {
					tempCoordinates = nil
					break
				} else {
					numberOfStepsLeft++
					tempCoordinates = append(tempCoordinates, currentPoint)
				}

			} else {
				tempCoordinates = append(tempCoordinates, currentPoint)
			}
		}
		coordinates = tempCoordinates
		ReserveSurroundingSlots(coordinates, occupiedFields)
	}
	return coordinates
}

func GeneratePlayerBattleships() [][]string {
	occupiedFields := make(map[string]bool)
	var shipCoordinates [][]string

	for i := 0; i < config.BATTLESHIP_COUNT; i++ {
		shipCoordinates = append(shipCoordinates, CreateBattleship(config.BATTLESHIP, occupiedFields))
	}

	for i := 0; i < config.SUBMARINE_COUNT; i++ {
		shipCoordinates = append(shipCoordinates, CreateBattleship(config.SUBMARINE, occupiedFields))
	}

	for i := 0; i < config.DESTROYER_COUNT; i++ {
		shipCoordinates = append(shipCoordinates, CreateBattleship(config.DESTROYER, occupiedFields))
	}

	for i := 0; i < config.PATROL_CRAFT_COUNT; i++ {
		shipCoordinates = append(shipCoordinates, CreateBattleship(config.PATROL_CRAFT, occupiedFields))
	}

	for i, ships := range shipCoordinates {
		shipCoordinates[i] = append([]string{strconv.Itoa(len(ships))}, ships...)
	}

	return shipCoordinates
}

func GenerateStringGrid(coordinates map[string]bool, isOpponentGrid bool) []string {
	currentPoint := ""
	var temp string
	var grid []string
	hit := config.SHIP
	miss := config.FIELD
	field := config.FIELD

	if isOpponentGrid {
		hit = config.HIT
		miss = config.MISSED
	}

	for row := 0; row < config.ROWS; row++ {
		temp = ""
		for column := 0; column < config.COLUMNS; column++ {
			currentPoint = fmt.Sprintf("%dx%s", availableRows[row], availableColumns[column])

			ship, check := coordinates[currentPoint]
			if ship {
				temp = temp + fmt.Sprintf("%s ", hit)
			} else if check {
				temp = temp + fmt.Sprintf("%s ", miss)
			} else {
				temp = temp + fmt.Sprintf("%s ", field)
			}

			if column == 9 {
				grid = append(grid, temp)
			}
		}
	}
	return grid
}

func ValidateCoordinates(coordinates []string) bool {
	//convert to runes

	startColumn := []rune(config.COLUMN_START)[0]
	endColumn := []rune(config.COLUMN_END)[0]
	checkup := true
	var spplitedCoordinates []string

	for i := 0; i < len(coordinates) && checkup; i++ {
		spplitedCoordinates = strings.Split(coordinates[0], config.COORDINATE_SEPARATOR)

		if len(spplitedCoordinates) != 2 {
			checkup = false
		}

		row, err := strconv.Atoi(spplitedCoordinates[0])
		if err != nil {
			checkup = false
		}

		column := []rune(spplitedCoordinates[1])[0]
		if row < config.ROW_START || row > config.ROW_END || column < startColumn || column > endColumn {
			checkup = false
		}

		if i == len(coordinates) {
			break
		}
	}
	return checkup
}

func FireShotsLogic(cordinates []string, gameId, playerId string) (responseData ResponseData, err error,
	responseStatus int, errorArg string) {

	game := &models.Game{}
	ships := &[]models.Battleship{}
	shot := &models.Shot{}
	userShots := &[]models.Shot{}
	battleship := &models.Battleship{}
	coordinatesMap := make(map[string]string)
	tempCoordinates := make(map[string]string)
	responseStruct := ResponseData{}
	insertShots := []models.Shot{}
	autopilot := &models.Autopilot{}

	battleship.AvailableShips(playerId, gameId, ships, config.DB)
	_, err, errArg, isFinished := game.FindGameById(gameId, playerId, config.DB)

	playerTurn := game.PlayerTurn
	won := game.Won
	opponentId := game.OpponentId
	numberOfLives := len(*ships)

	if len(cordinates) == 0 {
		cordinates = createRandomShots(numberOfLives)
	}

	areCordinatesValid := ValidateCoordinates(cordinates)

	if errArg != "" {
		return responseStruct, err, http.StatusNotFound, errArg
	}

	if !areCordinatesValid || len(cordinates) > numberOfLives {
		return responseStruct, nil, http.StatusBadRequest, ""
	}

	if isFinished {
		return responseStruct, nil, http.StatusMethodNotAllowed, ""
	}

	if fmt.Sprint(playerTurn) != playerId {
		return responseStruct, nil, http.StatusForbidden, ""
	}

	if fmt.Sprint(game.OpponentId) == playerId {
		opponentId = game.PlayerId
	}

	for _, coordinate := range cordinates {
		coordinatesMap[coordinate] = "MISS"
		tempCoordinates[coordinate] = "MISS"
	}

	//fetch previous shots
	_, err = shot.FindUserShots(fmt.Sprint(playerId), gameId, userShots, config.DB)

	if err != nil {
		return responseStruct, nil, http.StatusBadRequest, ""
	}

	for _, shot := range *userShots {
		//removed fired shots check
		if _, ok := coordinatesMap[shot.Coordinate]; ok {
			delete(tempCoordinates, shot.Coordinate)
		}
	}

	//fetch opponents battleship
	parsedOpponentId := strconv.FormatUint(uint64(opponentId), 10)
	_, err = battleship.AvailableShips(parsedOpponentId, gameId, ships, config.DB)

	//check shot status
	shotStatus(ships, tempCoordinates, coordinatesMap, gameId, playerId, won)

	responseStruct.Coordinates = coordinatesMap
	responseStruct.OpponentId = opponentId
	responseStruct.Won = won
	parsedPlayerId, err := strconv.Atoi(playerId)
	parsedGameId, err := strconv.Atoi(gameId)

	//adjust data for db
	for coordinate, value := range coordinatesMap {
		insertShots = append(insertShots, models.Shot{
			Coordinate: coordinate,
			Status:     value,
			PlayerId:   uint(parsedPlayerId),
			GameId:     uint(parsedGameId),
		})
	}

	_, err = shot.SaveShots(insertShots, config.DB)
	_, err = game.UpdateGameData(gameId, &models.Game{PlayerTurn: opponentId}, config.DB)

	if err != nil {
		return responseStruct, nil, http.StatusBadRequest, ""
	}

	//autopilot
	_, err = autopilot.CheckAutopilotData(opponentId, gameId, config.DB)
	if err != nil {
		tempArray := make([]string, 0)
		FireShotsLogic(tempArray, gameId, fmt.Sprint(opponentId))
	}

	return responseStruct, nil, http.StatusOK, ""
}

func shotStatus(ships *[]models.Battleship, tempCoordinates map[string]string, coordinatesMap map[string]string,
	gameId, playerId string, won string) {
	battleships := &models.Battleship{}
	game := &models.Game{}
	opponentLives := len(*ships)
	maxCordinatesCount := 4

	for _, battleship := range *ships {
		for i := 1; i <= maxCordinatesCount; i++ {
			currentField := fmt.Sprintf("Point%d", i)
			currentCoordinate := reflect.ValueOf(&battleship).Elem().FieldByName(currentField)
			currentCoordinateString := fmt.Sprint(currentCoordinate)

			if _, ok := tempCoordinates[currentCoordinateString]; ok && battleship.Lives == 1 {
				coordinatesMap[currentCoordinateString] = "KILL"
				delete(tempCoordinates, currentCoordinateString)
				battleships.UpdateBattleShipLives(battleship.Id, config.DB)
				opponentLives--

				if opponentLives == 0 {
					game.UpdateGameData(gameId, &models.Game{Won: playerId, Status: "FINISHED"}, config.DB)
					won = playerId
					break
				}
			} else if ok {
				coordinatesMap[currentCoordinateString] = "HIT"
				delete(tempCoordinates, currentCoordinateString)
				battleships.UpdateBattleShipLives(battleship.Id, config.DB)
				battleship.Lives -= 1
			}
		}

	}
}

func createRandomShots(limit int) []string {
	var tempCoordinates []string
	fmt.Println(tempCoordinates)
	for i := 0; i < limit; i++ {
		column := GetRandomValueInRange(65, 74)
		row := GetRandomValueInRange2(0, 10)
		coordinate := fmt.Sprintf("%dx%s", availableRows[row], string(column))
		tempCoordinates = append(tempCoordinates, coordinate)
	}

	return tempCoordinates
}
