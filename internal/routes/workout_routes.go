package routes

import (
	"github.com/amirintech/workout_go/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.App) *chi.Mux {
	router := chi.NewRouter()

	// health check route
	router.Get("/health", app.HealthCheck)

	// workout routes
	router.Get("/workout/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	router.Post("/workout", app.WorkoutHandler.HandlePostWorkout)
	router.Put("/workout/{id}", app.WorkoutHandler.HandlePutWorkout)
	router.Delete("/workout/{id}", app.WorkoutHandler.HandleDeleteWorkout)

	return router
}
