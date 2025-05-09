package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amirintech/workout_go/internal/api"
	"github.com/amirintech/workout_go/internal/store"
	"github.com/amirintech/workout_go/migrations"
)

type App struct {
	Logger         *log.Logger
	WorkoutHandler api.WorkoutHandler
	DB             *sql.DB
}

func New() (*App, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := store.Open()
	workoutHandler := api.NewWorkoutHandler(store.NewPostgresWorkoutStore(db))
	if err != nil {
		return nil, err
	}

	if err := store.MigrateFS(db, migrations.FS, "."); err != nil {
		panic(err)
	}

	app := &App{
		Logger:         logger,
		WorkoutHandler: *workoutHandler,
		DB:             db,
	}

	return app, nil
}

func (app *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is running\n")
}
