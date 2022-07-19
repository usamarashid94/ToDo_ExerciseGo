package router

import (
	"github.com/usamarashid94/ToDo_ExerciseGo/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/tasks", middleware.GetTasks).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/user", middleware.GetAllUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	// router.HandleFunc("/api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}
