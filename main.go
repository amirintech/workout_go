package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/amirintech/workout_go/internal/app"
	"github.com/amirintech/workout_go/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 3000, "sets server port")
	flag.Parse()

	app, err := app.New()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	router := routes.SetupRoutes(app)

	app.Logger.Printf("Running app on port %d...", port)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		app.Logger.Fatalf("Falied to start the server\nError: %v", err)
	}
}
