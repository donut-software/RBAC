package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
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
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.CREATE_ROLE_TASK)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	var req CreateRoleTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	// err = rb.svc.CreateRoleTask(r.Context(), req.TaskId, req.RoleId)
	err = rb.svc.CreateRoleTask(r.Context(), internal.RoleTasks{
		Task: internal.Tasks{
			Id: req.TaskId,
		},
		Role: internal.Roles{
			Id: req.RoleId,
		},
	})
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
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.GET_ROLE_TASK)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
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
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.UPDATE_ROLE_TASK)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	var req UpdateRoleTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	// err = rb.svc.UpdateRoleTask(r.Context(), req.TaskId, req.RoleId, req.Id)
	err = rb.svc.UpdateRoleTask(r.Context(), internal.RoleTasks{
		Id: req.Id,
		Task: internal.Tasks{
			Id: req.TaskId,
		},
		Role: internal.Roles{
			Id: req.RoleId,
		},
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating role", err)
		return
	}
	renderResponse(w,
		&RoleTaskResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}

type ListRoleTaskRequest struct {
	From int `json:"from"`
	Size int `json:"size"`
}
type ListRoleTaskResponse struct {
	RoleTask []RoleTask `json:"roletasks"`
	Total    int64      `json:"total"`
}

func (rb *RBACHandler) listRoleTask(w http.ResponseWriter, r *http.Request) {
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.LIST_ROLE_TASK)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	var req ListRoleTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	la, err := rb.svc.ListRoleTask(r.Context(), internal.ListArgs{
		From: &req.From,
		Size: &req.Size,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	roleTask := []RoleTask{}
	for _, value := range la.RoleTasks {
		//get task
		tk, err := rb.svc.Task(r.Context(), value.Task.Id)
		if err != nil {
			renderErrorResponse(r.Context(), w, "error getting task", err)
			return
		}
		task := Task{
			Id:        tk.Id,
			Task:      tk.Task,
			CreatedAt: tk.CreatedAt,
		}
		//get role
		rl, err := rb.svc.Role(r.Context(), value.Role.Id)
		if err != nil {
			renderErrorResponse(r.Context(), w, "error getting role", err)
			return
		}
		role := Role{
			Id:        rl.Id,
			Role:      rl.Role,
			CreatedAt: rl.CreatedAt,
		}
		roleTask = append(roleTask, RoleTask{
			Id:        value.Id,
			Task:      task,
			Role:      role,
			CreatedAt: value.CreatedAt,
		})
	}
	renderResponse(w, &ListRoleTaskResponse{
		RoleTask: roleTask,
		Total:    la.Total,
	}, http.StatusOK)
}

func (rb *RBACHandler) deleteRoleTask(w http.ResponseWriter, r *http.Request) {
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.DELETE_ROLE_TASK)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	roleTaskId := mux.Vars(r)["roleTaskId"]
	err = rb.svc.DeleteRoleTask(r.Context(), roleTaskId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error deleting roletask", err)
		return
	}
	renderResponse(w,
		&RoleTaskResponse{
			Message: "Deleted Successfully",
		}, http.StatusOK)
}
