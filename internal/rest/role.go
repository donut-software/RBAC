package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Role struct {
	Id        string    `json:"id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateRoleRequest struct {
	Role string `json:"role"`
}
type RoleResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) createRole(w http.ResponseWriter, r *http.Request) {
	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.CreateRole(r.Context(), req.Role)
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
	roleId := mux.Vars(r)["roleId"]
	role, err := rb.svc.Role(r.Context(), roleId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	renderResponse(w, &GetRoleResponse{
		Role: Role{
			Id:        role.Id,
			Role:      role.Role,
			CreatedAt: role.CreatedAt,
		},
	}, http.StatusOK)
}

type UpdateRoleRequest struct {
	RoleId string `json:"roleId"`
	Role   string `json:"role"`
}

func (rb *RBACHandler) updateRole(w http.ResponseWriter, r *http.Request) {
	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.UpdateRole(r.Context(), req.RoleId, req.Role)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating role", err)
		return
	}
	renderResponse(w,
		&RoleResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}
