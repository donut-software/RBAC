package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Task struct {
	Id        string    `json:"id"`
	Task      string    `json:"task"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateTaskRequest struct {
	Task string `json:"task"`
}
type TaskResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.CreateTask(r.Context(), req.Task)
	if err != nil {
		renderErrorResponse(r.Context(), w, "create task failed", err)
		return
	}
	renderResponse(w,
		&TaskResponse{
			Message: "Created Successfully",
		}, http.StatusCreated)
}

type GetTaskResponse struct {
	Task Task `json:"task"`
}

func (rb *RBACHandler) task(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["taskId"]
	task, err := rb.svc.Task(r.Context(), taskId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	renderResponse(w, &GetTaskResponse{
		Task: Task{
			Id:        task.Id,
			Task:      task.Task,
			CreatedAt: task.CreatedAt,
		},
	}, http.StatusOK)
}

type UpdateTaskRequest struct {
	TaskId string `json:"taskId"`
	Task   string `json:"task"`
}

func (rb *RBACHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.UpdateTask(r.Context(), req.TaskId, req.Task)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating task", err)
		return
	}
	renderResponse(w,
		&TaskResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}
