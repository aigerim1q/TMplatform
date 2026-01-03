package routes

import (
	"backend-go/controllers"

	"github.com/gorilla/mux"
)

func RegisterStageRoutes(r *mux.Router) {
	r.HandleFunc("/stages", controllers.CreateStage).Methods("POST")
	r.HandleFunc("/stages", controllers.GetStagesByProject).Methods("GET")
	r.HandleFunc("/stages/{id}", controllers.GetStageByID).Methods("GET")
	r.HandleFunc("/stages/{id}", controllers.UpdateStage).Methods("PUT")
	r.HandleFunc("/stages/{id}", controllers.DeleteStage).Methods("DELETE")
}
