package main

import (
	"battleships_backend/config"
	"battleships_backend/src/routes"
	"battleships_backend/src/socket"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	APIRoutes := routes.Api{}
	rand.Seed(time.Now().UnixNano())

	// load env file
	env := config.GetEnvironment()

	if err := godotenv.Load(".env." + env); err != nil {
		fmt.Println("Failed to load env")
		panic(err)
	}

	// Init the database
	config.Connect(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	// Init the routes
	APIRoutes.ServeRoutes()

	//init socket writer
	go socket.WriteSocketMessage()

	// Run the server
	fmt.Println("Connected To Database")
	fmt.Println("Server started port 8000")
	log.Fatal(http.ListenAndServe(":8000", APIRoutes.Router))
}
