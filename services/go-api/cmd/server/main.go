package main

import (
	"log"
	"net/http"

	"github.com/orcs-toolkit/orcs/services/go-api/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	log.Printf("go-api listening on %s", a.Config.ListenAddr)
	if err := http.ListenAndServe(a.Config.ListenAddr, a.Router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
