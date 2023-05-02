package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"CRUD_REST/db"
	"CRUD_REST/jwt"
	"CRUD_REST/models"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	userID, err := jwt.ValidateJWTToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	tasks, err := db.GetTasksFromDB(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tasks")
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) { // получение одной задачи по её ID
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get task")
		return
	}
	respondWithJSON(w, http.StatusOK, *task)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	taskID, err := db.CreateTaskInDB(newTask.UserID, newTask.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}
	newTask.ID = taskID
	respondWithJSON(w, http.StatusCreated, newTask)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	userID, err := jwt.ValidateJWTToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if task.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Task does not belong to user")
		return
	}

	err = db.UpdateTaskInDB(id, task)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	task.ID = id
	task.UserID = userID
	task.CompletedAt = time.Now()
	task.UpdatedAt = time.Now()

	respondWithJSON(w, http.StatusOK, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = db.DeleteTaskFromDB(taskID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete task from db")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Task deleted successfully"})
}
