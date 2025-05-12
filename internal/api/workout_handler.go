package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/amirintech/workout_go/internal/store"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(store store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{workoutStore: store}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		http.Error(w, "invalid workout ID", http.StatusBadRequest)
		return
	}

	fmt.Println("workoutID", workoutID)

	workout, err := wh.workoutStore.GetByID(int(workoutID))
	if err != nil {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(workout)
}

func (wh *WorkoutHandler) HandlePostWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		fmt.Println("error decoding workout", err)
		http.Error(w, "invalid workout data", http.StatusBadRequest)
		return
	}

	createdWorkout, err := wh.workoutStore.Create(&workout)
	if err != nil {
		fmt.Println("error creating workout", err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandlePutWorkout(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		http.Error(w, "invalid workout ID", http.StatusBadRequest)
		return
	}

	var updatedWorkout store.Workout
	if err := json.NewDecoder(r.Body).Decode(&updatedWorkout); err != nil {
		fmt.Println("error decoding workout", err)
		http.Error(w, "invalid workout data", http.StatusBadRequest)
		return
	}

	updatedWorkout.ID = int(workoutID)
	if err := wh.workoutStore.Update(&updatedWorkout); err != nil {
		fmt.Println("error updating workout", err)
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedWorkout)
}

func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deleting workout")
	id := chi.URLParam(r, "id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		http.Error(w, "invalid workout ID", http.StatusBadRequest)
		return
	}

	if err := wh.workoutStore.Delete(int(workoutID)); err != nil {
		fmt.Println("error deleting workout", err)
		http.Error(w, "failed to delete workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
