package main

import (
	"encoding/json"
	"fmt"
	"igabir98/simpleTODO/engine"
	"igabir98/simpleTODO/models"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	db *engine.BoltDB
}

var app App

func main() {
	db, err := engine.NewBoltDB()

	if err != nil {
		fmt.Println(err)
	}

	app.db = db

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", welcome)
	r.Get("/tasks", getAll)
	r.Get("/tasks/{taskID}", getTask)
	r.Post("/tasks", create)
	r.Put("/tasks", update)
	r.Delete("/tasks/{taskID}", deleteTask)

	http.ListenAndServe(":3000", r)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func create(w http.ResponseWriter, req *http.Request) {
	var task models.Task
	json.NewDecoder(req.Body).Decode(&task)
	w.Header().Set("Content-Type", "application/json")

	_, err := app.db.CreateTask(&task)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(task)
}

func update(w http.ResponseWriter, req *http.Request) {
	var task models.Task
	json.NewDecoder(req.Body).Decode(&task)

	_, err := app.db.UpdateTask(&task)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(task)
}

func getAll(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)

	tasks, _ := app.db.GetAllTasks()

	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, req *http.Request) {
	taskID := chi.URLParam(req, "taskID")

	w.Header().Set("Content-Type", "application/json")

	tasks, err := app.db.GetTask(taskID)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusFound)

	json.NewEncoder(w).Encode(tasks)
}

func deleteTask(w http.ResponseWriter, req *http.Request) {
	taskID := chi.URLParam(req, "taskID")

	w.Header().Set("Content-Type", "application/json")

	err := app.db.DeleteTask(taskID)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}
