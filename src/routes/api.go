package routes

import (
	"battleships_backend/src/controllers"
	"battleships_backend/src/middlewares"
	socket "battleships_backend/src/socket"

	"github.com/gorilla/mux"
)

type Api struct {
	Router *mux.Router
}

// ServeRoutes handle the public routes
func (api *Api) ServeRoutes() {
	api.Router = mux.NewRouter()

	//init websocket
	api.Router.HandleFunc("/ws", socket.SocketEndPoint)

	// Route List
	UserRouter := api.Router.PathPrefix("/player").Subrouter()

	// Middleware
	UserRouter.Use(middlewares.SetContentTypeHeader)

	// Routes
	UserRouter.HandleFunc("/{player_id}/game/list", controllers.GetAllGamesByUser).Methods("GET")
	UserRouter.HandleFunc("", controllers.Register).Methods("POST")
	UserRouter.HandleFunc("/list", controllers.GetAvailableUser).Methods("GET")
	UserRouter.HandleFunc("/{id}", controllers.GetUserByID).Methods("GET")
	UserRouter.HandleFunc("/{opponent_id}/game", controllers.CreateGame).Methods("POST")
	UserRouter.HandleFunc("/{player_id}/game/{game_id}", controllers.ViewGameStatus).Methods("GET")
	UserRouter.HandleFunc("/{player_id}/game/{game_id}", controllers.FireShots).Methods("PUT")
	UserRouter.HandleFunc("/{player_id}/game/{game_id}/autopilot", controllers.TurnAutopilot).Methods("PUT")
}
