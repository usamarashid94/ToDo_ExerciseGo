package router

import (
	"github.com/usamarashid94/ToDo_ExerciseGo/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/tasks", middleware.GetTasks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", middleware.AddTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deletetask/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	// router.HandleFunc("/api/user", middleware.GetAllUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	//

	return router
}
