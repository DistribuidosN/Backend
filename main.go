package main

import (
	"Backend/routes"
	"log"
)

func main() {
	r := routes.SetupRoutes()
	log.Println("Server starting on port 50021...")
	if err := r.Run(":50021"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
