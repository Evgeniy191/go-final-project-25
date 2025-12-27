package server

import (
	"fmt"
	"go-final-project-25/pkg/api"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Start() error {
	api.Init()

	port := 7540

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Server started on http://localhost:%d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
