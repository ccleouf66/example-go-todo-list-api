package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/codenito/example-go-todo-list-api/pkg/store"
	"github.com/codenito/example-go-todo-list-api/pkg/types"
)

type TaskHandler struct {
	Store *store.MongoStore
}

func returnError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, fmt.Sprintf("{\"status\": \"error\", \"message\": \"%s\"}", err))
	log.Printf("%s\n", err)
}

func (h *TaskHandler) ServeHTTP(r chi.Router) {
	r.Get("/", h.getTasks)
	r.Post("/", h.createTask)
	r.Delete("/", h.deleteTask)
}

func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Store.GetTasks(r.Context())
	if err != nil {
		returnError(w, err)
		return
	}

	// Convert tasks to json
	b, err := json.Marshal(tasks)
	if err != nil {
		returnError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(b))
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var taskToCreate types.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&taskToCreate)
	if err != nil {
		returnError(w, err)
		return
	}
	if taskToCreate.Date.IsZero() {
		taskToCreate.Date = time.Now()
	}

	taskInDb, err := h.Store.CreateTask(r.Context(), taskToCreate)
	if err != nil {
		returnError(w, err)
		return
	}

	// Convert tasks to json
	b, err := json.Marshal(taskInDb)
	if err != nil {
		returnError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, string(b))

}
func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	var taskToDelete types.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&taskToDelete)
	if err != nil {
		returnError(w, err)
		return
	}

	err = h.Store.DeleteTask(r.Context(), taskToDelete)
	if err != nil {
		returnError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
