package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go-final-project-25/pkg/db"
	"go-final-project-25/pkg/nextdate"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{"error": "Не указан идентификатор"})
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]any{"error": "Задача не найдена"})
		return
	}

	// Возвращаем задачу
	writeJSON(w, task)
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{"error": "Ошибка десериализации JSON"})
		return
	}

	// Проверяем ID
	if task.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{"error": "Не указан идентификатор"})
		return
	}

	// Проверяем заголовок
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{"error": "Не указан заголовок задачи"})
		return
	}

	_, err := time.Parse(nextdate.DateFormat, task.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{"error": "Некорректный формат даты"})
		return
	}

	if task.Repeat != "" {
		_, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, map[string]any{"error": err.Error()})
			return
		}
	}

	err = db.UpdateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]any{"error": "Задача не найдена"})
		return
	}

	writeJSON(w, map[string]any{})
}
