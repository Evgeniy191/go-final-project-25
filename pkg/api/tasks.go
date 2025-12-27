package api

import (
	"net/http"

	"go-final-project-25/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// только GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]any{"error": "method not allowed"})
		return
	}

	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search == "" {
		tasks, err = db.Tasks(50)
	} else {
	
		tasks, err = db.SearchTasks(search)
	}

	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
