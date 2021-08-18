package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Role struct {
	Id        string    `json:"id"`
	Role      string    `json:"role"`
	Task      []Task    `json:"tasks"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateRoleRequest struct {
	Role string `json:"role"`
}
type RoleResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) createRole(w http.ResponseWriter, r *http.Request) {
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.CREATE_ROLE)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	_, err = rb.svc.CreateRole(r.Context(), req.Role)
	if err != nil {
		renderErrorResponse(r.Context(), w, "create role failed", err)
		return
	}
	renderResponse(w,
		&RoleResponse{
			Message: "Created Successfully",
		}, http.StatusCreated)
}

type GetRoleResponse struct {
	Role Role `json:"role"`
}

func (rb *RBACHandler) role(w http.ResponseWriter, r *http.Request) {
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.GET_ROLE)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	roleId := mux.Vars(r)["roleId"]
	// role, err := rb.svc.Role(r.Context(), roleId)
	// if err != nil {
	// 	renderErrorResponse(r.Context(), w, "error getting the role", err)
	// 	return
	// }
	//roletaskbyrole
	rt, err := rb.svc.RoleTaskByRole(r.Context(), roleId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting task by role", err)
		return
	}
	var tasks []Task
	for _, value := range rt.Tasks {
		tasks = append(tasks, Task{
			Id:        value.Id,
			Task:      value.Task,
			CreatedAt: value.CreatedAt,
		})
	}
	renderResponse(w, &GetRoleResponse{
		Role: Role{
			Id:        rt.Role.Id,
			Role:      rt.Role.Role,
			Task:      tasks,
			CreatedAt: rt.Role.CreatedAt,
		},
	}, http.StatusOK)
}

type UpdateRoleRequest struct {
	RoleId string `json:"roleId"`
	Role   string `json:"role"`
}

func (rb *RBACHandler) updateRole(w http.ResponseWriter, r *http.Request) {
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.UPDATE_ROLE)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	// err = rb.svc.UpdateRole(r.Context(), req.RoleId, req.Role)
	err = rb.svc.UpdateRole(r.Context(), internal.Roles{
		Id:   req.RoleId,
		Role: req.Role,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating role", err)
		return
	}
	renderResponse(w,
		&RoleResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}

type ListRoleResponse struct {
	Roles []Role `json:"roles"`
	Total int64  `json:"total"`
}

func (rb *RBACHandler) listrole(w http.ResponseWriter, r *http.Request) {
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.LIST_ROLE)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	var from int
	var size int
	v := r.URL.Query()
	from, err = strconv.Atoi(v.Get("from"))
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid param from", err)
		return
	}
	size, err = strconv.Atoi(v.Get("size"))
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid param size", err)
		return
	}
	la, err := rb.svc.ListRole(r.Context(), internal.ListArgs{
		From: &from,
		Size: &size,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	roles := []Role{}
	for _, value := range la.Roles {
		rt, err := rb.svc.RoleTaskByRole(r.Context(), value.Id)
		if err != nil {
			renderErrorResponse(r.Context(), w, "error getting task by role", err)
			return
		}
		var tasks []Task
		for _, value := range rt.Tasks {
			tasks = append(tasks, Task{
				Id:        value.Id,
				Task:      value.Task,
				CreatedAt: value.CreatedAt,
			})
		}
		acc := Role{
			Id:        value.Id,
			Role:      value.Role,
			Task:      tasks,
			CreatedAt: value.CreatedAt,
		}
		roles = append(roles, acc)
	}
	renderResponse(w, &ListRoleResponse{
		Roles: roles,
		Total: la.Total,
	}, http.StatusOK)
}

type AccountRoleByRole struct {
	Role    Role      `json:"role"`
	Account []Account `json:"accounts"`
}

func (rb *RBACHandler) getAccountRoleByRole(w http.ResponseWriter, r *http.Request) {
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.GET_ROLE)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	id := mux.Vars(r)["roleId"]
	la, err := rb.svc.AccountRoleByRole(r.Context(), id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	role := Role{
		Id:        la.Role.Id,
		Role:      la.Role.Role,
		CreatedAt: la.Role.CreatedAt,
	}
	account := []Account{}
	for _, value := range la.Account {
		prof := Profile{
			Id:                value.Profile.Id,
			ProfilePicture:    value.Profile.Profile_Picture,
			ProfileBackground: value.Profile.Profile_Background,
			FirstName:         value.Profile.First_Name,
			LastName:          value.Profile.Last_Name,
			Mobile:            value.Profile.Mobile,
			Email:             value.Profile.Email,
			CreatedAt:         value.CreatedAt,
		}
		acc := Account{
			Id:        value.Id,
			Username:  value.UserName,
			Profile:   prof,
			CreatedAt: value.CreatedAt,
		}
		account = append(account, acc)
	}
	renderResponse(w, &AccountRoleByRole{
		Role:    role,
		Account: account,
	}, http.StatusOK)
}

func (rb *RBACHandler) deleteRole(w http.ResponseWriter, r *http.Request) {
	authusername := r.Header.Get("username")
	allowed, err := rb.svc.IsAllowed(r.Context(), authusername, internal.DELETE_ROLE)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting user tasks", err)
		return
	}
	if !allowed {
		renderErrorResponse(r.Context(), w, "user is not allowed", err)
		return
	}
	roleId := mux.Vars(r)["roleId"]
	err = rb.svc.DeleteRole(r.Context(), roleId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error deleting role", err)
		return
	}
	renderResponse(w,
		&RoleResponse{
			Message: "Deleted Successfully",
		}, http.StatusOK)
}
