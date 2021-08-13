package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type RoleTask struct {
	Id        string    `json:"id"`
	Task      Task      `json:"task"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateRoleTaskRequest struct {
	TaskId string `json:"taskId"`
	RoleId string `json:"roleId"`
}
type RoleTaskResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) createRoleTask(w http.ResponseWriter, r *http.Request) {
	var req CreateRoleTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.CreateRoleTask(r.Context(), req.TaskId, req.RoleId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "create accountrole failed", err)
		return
	}
	renderResponse(w,
		&RoleTaskResponse{
			Message: "Created Successfully",
		}, http.StatusCreated)
}

type GetRoleTaskResponse struct {
	RoleTask RoleTask `json:"roleTask"`
}

func (rb *RBACHandler) roleTask(w http.ResponseWriter, r *http.Request) {
	roleTaskId := mux.Vars(r)["roleTaskId"]
	roleTask, err := rb.svc.RoleTask(r.Context(), roleTaskId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	rl, err := rb.svc.Role(r.Context(), roleTask.Role.Id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	ts, err := rb.svc.Task(r.Context(), roleTask.Task.Id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}

	role := Role{
		Id:        rl.Id,
		Role:      rl.Role,
		CreatedAt: rl.CreatedAt,
	}
	task := Task{
		Id:        ts.Id,
		Task:      ts.Task,
		CreatedAt: ts.CreatedAt,
	}
	renderResponse(w, &GetRoleTaskResponse{
		RoleTask: RoleTask{
			Id:        roleTask.Id,
			Task:      task,
			Role:      role,
			CreatedAt: roleTask.CreatedAt,
		},
	}, http.StatusOK)
}

type UpdateRoleTaskRequest struct {
	TaskId string `json:"taskid"`
	RoleId string `json:"roleId"`
	Id     string `json:"id"`
}

func (rb *RBACHandler) updateRoleTask(w http.ResponseWriter, r *http.Request) {
	var req UpdateRoleTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.UpdateRoleTask(r.Context(), req.TaskId, req.RoleId, req.Id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating role", err)
		return
	}
	renderResponse(w,
		&RoleTaskResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}
