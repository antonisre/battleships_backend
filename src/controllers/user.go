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

	"github.com/gorilla/mux"
)

// Register a new user
func Register(response http.ResponseWriter, request *http.Request) {
	user := &models.User{}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	// Validate the user
	err = user.ValidateRegister(config.DB)
	if err != nil {
		lib.Error(response, http.StatusUnprocessableEntity, err.Error(), "")
		return
	}

	// Check the user
	checkUser, _ := user.GetUserByEmail(config.DB)
	if checkUser != nil {
		lib.Error(response, http.StatusConflict, "error.username-already-taken", user.Email)
		return
	}

	userData, err := user.Register(config.DB)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error(), "")
		return
	}

	data := map[string]string{
		"name":  userData.Name,
		"email": userData.Email,
	}

	response.Header().Set("Location", fmt.Sprintf("/player/player-%d", user.ID))
	lib.Success(response, http.StatusCreated, data)
	return
}

// Find user by id
func GetUserByID(response http.ResponseWriter, request *http.Request) {
	user := &models.UserJSON{}
	params := mux.Vars(request)

	userData, err := user.GetUser(params["id"], config.DB)
	if err != nil {
		switch err.Error() {
		case "record not found":
			lib.Error(response, http.StatusNotFound, "", "")
			return
		default:
			lib.Error(response, http.StatusBadRequest, err.Error(), "")
			return
		}
	}

	data := map[string]string{
		"name":  userData.Name,
		"email": userData.Email,
	}

	lib.Success(response, http.StatusOK, data)
	return
}

func GetAvailableUser(response http.ResponseWriter, request *http.Request) {
	user := &models.UserJSON{}
	queryParams := request.URL.Query()
	var page int = 1
	var limitPerPage int = 50

	if pageInput, ok := (queryParams["page"]); ok {
		if tempPage, err := strconv.Atoi(pageInput[0]); err == nil {
			page = tempPage
		}
	}

	offset := limitPerPage * (page - 1)
	userData, err := user.GetUsers(offset, limitPerPage, config.DB)
	if err != nil {
		panic(err.Error())
	}

	players := map[string][]models.UserJSON{
		"players": *userData,
	}

	lib.Success(response, http.StatusOK, players)
	return
}
