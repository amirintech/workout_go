package routes

import (
	"github.com/amirintech/workout_go/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)
	router.Get("/workout/{id}", app.WorkoutHandler.HandleGetWorkoutByID)

	router.Post("/workout", app.WorkoutHandler.HandlePostWorkout)

	return router
}
