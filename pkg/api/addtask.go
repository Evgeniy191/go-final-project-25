package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go-final-project-25/pkg/db"
	"go-final-project-25/pkg/nextdate"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]any{"error": "title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]any{"id": id})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(nextdate.DateFormat)
		return nil
	}

	t, err := time.Parse(nextdate.DateFormat, task.Date)
	if err != nil {
		return err
	}

	var next string
	if task.Repeat != "" {
		next, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if dateOnly.Before(nowOnly) {
		if task.Repeat == "" {
			task.Date = now.Format(nextdate.DateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}
