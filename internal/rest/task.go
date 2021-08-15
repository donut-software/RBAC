package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
	"time"

	"github.com/gorilla/mux"
)

type Task struct {
	Id         string       `json:"id"`
	Task       string       `json:"task"`
	HelpText   HelpText     `json:"helptext"`
	Menus      []Menu       `json:"menus"`
	Navigation []Navigation `json:"navigation"`
	CreatedAt  time.Time    `json:"created_at"`
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
	_, err := rb.svc.CreateTask(r.Context(), req.Task)
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
			Id:   task.Id,
			Task: task.Task,
			HelpText: HelpText{
				Id:        task.HelpText.Id,
				HelpText:  task.HelpText.HelpText,
				TaskId:    task.HelpText.Task_id,
				CreatedAt: task.HelpText.CreatedAt,
			},
			Menus:      convertInternalMenuList(task.Menu),
			Navigation: convertInternalNavigationList(task.Navigation),
			CreatedAt:  task.CreatedAt,
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

type ListTaskRequest struct {
	From int `json:"from"`
	Size int `json:"size"`
}
type ListTaskResponse struct {
	Tasks []Task `json:"tasks"`
	Total int64  `json:"total"`
}

func (rb *RBACHandler) listtask(w http.ResponseWriter, r *http.Request) {
	var req ListTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	la, err := rb.svc.ListTask(r.Context(), internal.ListArgs{
		From: &req.From,
		Size: &req.Size,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	tasks := []Task{}
	for _, value := range la.Task {
		acc := Task{
			Id:   value.Id,
			Task: value.Task,
			HelpText: HelpText{
				Id:        value.HelpText.Id,
				HelpText:  value.HelpText.HelpText,
				TaskId:    value.HelpText.Task_id,
				CreatedAt: value.HelpText.CreatedAt,
			},
			Menus:      convertInternalMenuList(value.Menu),
			Navigation: convertInternalNavigationList(value.Navigation),
			CreatedAt:  value.CreatedAt,
		}
		tasks = append(tasks, acc)
	}
	renderResponse(w, &ListTaskResponse{
		Tasks: tasks,
		Total: la.Total,
	}, http.StatusOK)
}

func (rb *RBACHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["taskId"]
	err := rb.svc.DeleteTask(r.Context(), taskId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error deleting task", err)
		return
	}
	renderResponse(w,
		&TaskResponse{
			Message: "Deleted Successfully",
		}, http.StatusOK)
}
