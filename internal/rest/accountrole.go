package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
	"time"

	"github.com/gorilla/mux"
)

type AccountRole struct {
	Id        string    `json:"id"`
	Account   Account   `json:"account"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateAccountRoleRequest struct {
	AccountId string `json:"account_id"`
	RoleId    string `json:"role_id"`
}
type AccountRoleResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) createAccountRole(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.CreateAccountRole(r.Context(), req.AccountId, req.RoleId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "create accountrole failed", err)
		return
	}
	renderResponse(w,
		&AccountRoleResponse{
			Message: "Created Successfully",
		}, http.StatusCreated)
}

type GetAccountRoleResponse struct {
	AccountRole AccountRole `json:"accountRole"`
}

func (rb *RBACHandler) accountRole(w http.ResponseWriter, r *http.Request) {
	accountRoleId := mux.Vars(r)["accountRoleId"]
	accountRole, err := rb.svc.AccountRole(r.Context(), accountRoleId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	acc, err := rb.svc.AccountByID(r.Context(), accountRole.Account.Id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	rl, err := rb.svc.Role(r.Context(), accountRole.Role.Id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}

	profile := Profile{
		Id:                acc.Profile.Id,
		ProfileBackground: acc.Profile.Profile_Background,
		ProfilePicture:    acc.Profile.Profile_Picture,
		FirstName:         acc.Profile.First_Name,
		LastName:          acc.Profile.Last_Name,
		Mobile:            acc.Profile.Mobile,
		Email:             acc.Profile.Email,
		CreatedAt:         acc.CreatedAt,
	}
	account := Account{
		Id:        acc.Id,
		Username:  acc.UserName,
		Profile:   profile,
		CreatedAt: acc.CreatedAt,
	}

	role := Role{
		Id:        rl.Id,
		Role:      rl.Role,
		CreatedAt: rl.CreatedAt,
	}
	renderResponse(w, &GetAccountRoleResponse{
		AccountRole: AccountRole{
			Id:        accountRole.Id,
			Account:   account,
			Role:      role,
			CreatedAt: accountRole.CreatedAt,
		},
	}, http.StatusOK)
}

type UpdateAccountRoleRequest struct {
	AccountId string `json:"accountId"`
	RoleId    string `json:"roleId"`
	Id        string `json:"id"`
}

func (rb *RBACHandler) updateAccountRole(w http.ResponseWriter, r *http.Request) {
	var req UpdateAccountRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.UpdateAccountRole(r.Context(), req.AccountId, req.RoleId, req.Id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating role", err)
		return
	}
	renderResponse(w,
		&RoleResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}

type ListAccountRoleRequest struct {
	From int `json:"from"`
	Size int `json:"size"`
}
type ListAccountRoleResponse struct {
	AccoutRoles []AccountRole `json:"accountRoles"`
	Total       int64         `json:"total"`
}

func (rb *RBACHandler) listAccountRole(w http.ResponseWriter, r *http.Request) {
	var req ListRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	la, err := rb.svc.ListAccountRole(r.Context(), internal.ListArgs{
		From: &req.From,
		Size: &req.Size,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	acRoles := []AccountRole{}
	for _, value := range la.AccountRoles {
		//get profile
		acc, err := rb.svc.Account(r.Context(), value.Account.UserName)
		if err != nil {
			renderErrorResponse(r.Context(), w, "error getting account", err)
			return
		}
		profile := Profile{
			Id:                acc.Profile.Id,
			ProfileBackground: acc.Profile.Profile_Background,
			ProfilePicture:    acc.Profile.Profile_Picture,
			FirstName:         acc.Profile.First_Name,
			LastName:          acc.Profile.Last_Name,
			Mobile:            acc.Profile.Mobile,
			Email:             acc.Profile.Email,
			CreatedAt:         acc.CreatedAt,
		}
		account := Account{
			Id:        acc.Id,
			Username:  acc.UserName,
			Profile:   profile,
			CreatedAt: acc.CreatedAt,
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
		acRoles = append(acRoles, AccountRole{
			Id:        value.Id,
			Account:   account,
			Role:      role,
			CreatedAt: value.CreatedAt,
		})
	}
	renderResponse(w, &ListAccountRoleResponse{
		AccoutRoles: acRoles,
		Total:       la.Total,
	}, http.StatusOK)
}

func (rb *RBACHandler) deleteAccountRole(w http.ResponseWriter, r *http.Request) {
	accountRoleId := mux.Vars(r)["accountRoleId"]
	err := rb.svc.DeleteAccountRole(r.Context(), accountRoleId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error deleting accountRole", err)
		return
	}
	renderResponse(w,
		&AccountRoleResponse{
			Message: "Deleted Successfully",
		}, http.StatusOK)
}
