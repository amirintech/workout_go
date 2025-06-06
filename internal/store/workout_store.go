package store

import (
	"database/sql"
	"fmt"
)

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type WorkoutStore interface {
	Create(workout *Workout) (*Workout, error)
	GetByID(id int) (*Workout, error)
	Update(workout *Workout) error
	Delete(id int) error
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

func (pws *PostgresWorkoutStore) Create(workout *Workout) (*Workout, error) {
	tx, err := pws.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO workouts (title, description, duration_minutes, calories_burned)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	err = tx.QueryRow(
		query,
		workout.Title,
		workout.Description,
		workout.DurationMinutes,
		workout.CaloriesBurned,
	).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, weight, duration_seconds, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
		`
		err = tx.QueryRow(
			query,
			workout.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.Weight,
			entry.DurationSeconds,
			entry.Notes,
			entry.OrderIndex,
		).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return workout, nil
}

func (pws *PostgresWorkoutStore) GetByID(id int) (*Workout, error) {
	query := `
	SELECT id, title, description, duration_minutes, calories_burned
	FROM workouts
	WHERE id = $1;
	`
	var workout Workout
	err := pws.db.QueryRow(query, id).Scan(
		&workout.ID,
		&workout.Title,
		&workout.Description,
		&workout.DurationMinutes,
		&workout.CaloriesBurned,
	)

	if err != nil {
		return nil, err
	}

	workout.Entries, err = pws.getEntriesForWorkout(id)
	if err != nil {
		return nil, err
	}

	return &workout, nil
}

func (pws *PostgresWorkoutStore) getEntriesForWorkout(id int) ([]WorkoutEntry, error) {
	query := `
	SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
	FROM workout_entries
	WHERE workout_id = $1
	ORDER BY order_index;
	`
	rows, err := pws.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []WorkoutEntry{}
	for rows.Next() {
		var entry WorkoutEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (pws *PostgresWorkoutStore) Update(workout *Workout) error {
	tx, err := pws.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	UPDATE workouts
	SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
	WHERE id = $5;
	`
	res, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("workout not found")
	}

	for _, entry := range workout.Entries {
		query := `
		UPDATE workout_entries
		SET exercise_name = $1, sets = $2, reps = $3, duration_seconds = $4, weight = $5, notes = $6, order_index = $7
		WHERE id = $8;
		`
		res, err := tx.Exec(query, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex, entry.ID)
		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("workout entry not found")
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (pws *PostgresWorkoutStore) Delete(id int) error {
	query := `
	DELETE FROM workouts WHERE id = $1;
	`
	res, err := pws.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("workout not found")
	}

	return nil
}
