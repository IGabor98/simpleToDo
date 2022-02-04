package main

import (
	"encoding/json"
	"igabir98/simpleTODO/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var tasks []models.Task

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", welcome)
	r.Post("/tasks", create)
	r.Get("/tasks", getAll)

	http.ListenAndServe(":3000", r)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func create(w http.ResponseWriter, req *http.Request) {
	var task models.Task
	json.NewDecoder(req.Body).Decode(&task)

	tasks = append(tasks, task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func getAll(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(tasks)
}
