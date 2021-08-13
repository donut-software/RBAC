package rest

import (
	"encoding/json"
	"net/http"
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
		ProfileBackground: acc.Profile.Profile_Background,
		ProfilePicture:    acc.Profile.Profile_Picture,
		FirstName:         acc.Profile.First_Name,
		LastName:          acc.Profile.Last_Name,
		Mobile:            acc.Profile.Mobile,
		Email:             acc.Profile.Email,
		CreatedAt:         acc.CreatedAt,
	}
	account := Account{
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
