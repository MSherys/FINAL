package nextdate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go1f/pkg/constants"
	dbase "go1f/pkg/db"
)

var task dbase.Task

func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	type errorjson struct {
		Error string `json:"error"`
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if !validateAndAdjustTask(&task, w) {
		return
	}

	id, err := dbase.AddTask(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(errorjson{Error: "ошибка AddTask"}); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]interface{}{"id": id}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	return
}

func validateAndAdjustTask(task *dbase.Task, w http.ResponseWriter) bool {
	type errorjson struct {
		Error string `json:"error"`
	}

	if len(task.Title) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errorjson{Error: "Не указан заголовок задачи"}); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return false
	}

	now := time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format(constants.DateFormat)
	} else {
		t, err := time.Parse(constants.DateFormat, task.Date)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(errorjson{Error: fmt.Sprintf("дата представлена в формате, отличном от %s", constants.DateFormat)}); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return false
		}

		var nextd string
		if len(task.Repeat) > 0 {
			nextd, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(errorjson{Error: "ошибка NextDate"}); err != nil {
					http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
				}
				return false
			}
		}

		if afterNow(now, t) {
			if len(task.Repeat) == 0 {
				task.Date = now.Format(constants.DateFormat)
			} else {
				task.Date = nextd
			}
		}
	}
	return true
}
