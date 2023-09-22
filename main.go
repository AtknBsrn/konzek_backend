package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/validator.v2"
)

var db *sql.DB

// TaskQueue is a buffered channel for tasks
var TaskQueue chan Task

// ErrQueue is a buffered channel for errors
var ErrQueue chan error

const NumberOfWorkers = 5

var wg sync.WaitGroup

func main() {
	create_server()
}

func create_server() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Loading Server Database Connection...")
	create_database()
	log.Println("Creating Worker Pool...")
	createWorkerPool()

	r := mux.NewRouter()
	r.HandleFunc("/tasks", withAuth(CreateTaskRequest)).Methods("POST")
	r.HandleFunc("/tasks/{id}", withAuth(GetTaskRequest)).Methods("GET")
	r.HandleFunc("/tasks/{id}", withAuth(UpdateTaskRequest)).Methods("PUT")
	r.HandleFunc("/tasks/{id}", withAuth(DeleteTaskRequest)).Methods("DELETE")

	log.Println("Serving HTTP :8080...")
	http_err := http.ListenAndServe(":8080", r)
	if http_err != nil {
		log.Println("HTTP SERVER FAILED", http_err)
	}

	close(TaskQueue)
	wg.Wait()

	close(ErrQueue)
	for err := range ErrQueue {
		log.Printf("Queued error: %v", err)
	}
}

// createTaskRequest is an HTTP handler function that handles the creation
// of a new task. It expects a Task object in the HTTP request body.
func CreateTaskRequest(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		response := JsonResponse{"Error", "Invalid JSON provided"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if errs := validator.Validate(task); errs != nil {
		response := JsonResponse{"Error", errs.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	select {
	case TaskQueue <- task:
		response := JsonResponse{"Accepted", "Task creation in progress"}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	default:
		response := JsonResponse{"Error", "Task queue is full. Please try again later."}
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
	}
}

// getTaskRequest is an HTTP handler function that handles fetching a task by ID.
// It expects a task ID in the request URL.
func GetTaskRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response := JsonResponse{"Error", "Invalid ID provided"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	task, err := getTask(id)
	if err != nil {
		response := JsonResponse{"Error", err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// updateTaskRequest is an HTTP handler function that handles updating a task by ID.
// It expects a task ID in the request URL and a Task object in the request body
func UpdateTaskRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response := JsonResponse{"Error", "Invalid ID provided"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var t Task
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		response := JsonResponse{"Error", "Invalid JSON provided"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	t.ID = id
	err = updateTask(t)
	if err != nil {
		response := JsonResponse{"Error", err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := JsonResponse{"Success", "Task updated successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// deleteTaskRequest is an HTTP handler function that handles deleting a task by ID.
// It expects a task ID in the request URL.
func DeleteTaskRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response := JsonResponse{"Error", "Invalid ID provided"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = deleteTask(id)
	if err != nil {
		response := JsonResponse{"Error", err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := JsonResponse{"Success", "Task deleted successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
