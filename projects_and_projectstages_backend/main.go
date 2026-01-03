package main

import (
	"fmt"
	"log"
	"net/http"

	"backend-go/db"
	"backend-go/routes"

	"github.com/gorilla/mux"
)

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	routes.RegisterProjectRoutes(r)
	routes.RegisterStageRoutes(r)

	fmt.Println("Server running on port 3001")
	log.Fatal(http.ListenAndServe(":3001", r))
}
