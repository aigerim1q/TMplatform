package routes

import (
	"backend-go/controllers"

	"github.com/gorilla/mux"
)

func RegisterProjectRoutes(r *mux.Router) {
	r.HandleFunc("/projects", controllers.CreateProject).Methods("POST")
	r.HandleFunc("/projects", controllers.GetMyProjects).Methods("GET")
	r.HandleFunc("/projects/{id}", controllers.GetProjectByID).Methods("GET")
	r.HandleFunc("/projects/{id}", controllers.UpdateProject).Methods("PUT")
	r.HandleFunc("/projects/{id}", controllers.DeleteProject).Methods("DELETE")
}
