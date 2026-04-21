package main

import (
	"Backend/config"
	"Backend/routes"
	"log"
)

func main() {
	cfg := config.Load()
	r := routes.SetupRoutes(cfg)
	log.Printf("Server starting on %s and proxying SOAP to %s", cfg.RESTBindAddr, cfg.ServerAppSOAPBase)
	if err := r.Run(cfg.RESTBindAddr); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
