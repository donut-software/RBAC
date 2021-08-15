package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
	"time"

	"github.com/gorilla/mux"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Message string `json:"message"`
}

func (a *RBACHandler) logout(w http.ResponseWriter, r *http.Request) {
	addCookie(w, "token", "")
	renderResponse(w,
		&LoginResponse{
			Message: "Logout Succesfully",
		}, http.StatusCreated)
}

func (a *RBACHandler) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	defer r.Body.Close()
	err := a.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		renderErrorResponse(r.Context(), w, "login failed", err)
		return
	}
	cookie, err := a.svc.CreateToken(req.Username)
	if err != nil {
		renderErrorResponse(r.Context(), w, "token creation failed", err)
		return
	}
	addCookie(w, "token", cookie)
	renderResponse(w,
		&LoginResponse{
			Message: "Login Succesfully",
		}, http.StatusCreated)
}

func (a *RBACHandler) me(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	account, err := a.svc.Account(r.Context(), username)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	profile := Profile{
		Id:                account.Profile.Id,
		ProfileBackground: account.Profile.Profile_Background,
		ProfilePicture:    account.Profile.Profile_Picture,
		FirstName:         account.Profile.First_Name,
		LastName:          account.Profile.Last_Name,
		Mobile:            account.Profile.Mobile,
		Email:             account.Profile.Email,
		CreatedAt:         account.CreatedAt,
	}
	renderResponse(w, &ReadAccountResponse{
		Account: Account{
			Id:        account.Id,
			Username:  account.UserName,
			Profile:   profile,
			CreatedAt: account.CreatedAt,
		},
	}, http.StatusOK)
}

type Profile struct {
	Id                string    `json:"id"`
	ProfilePicture    string    `json:"profile_picture"`
	ProfileBackground string    `json:"profile_background"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Mobile            string    `json:"mobile"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
}

type Account struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Profile   Profile   `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	ProfilePicure     string `json:"profile_picture"`
	ProfileBackground string `json:"profile_background"`
	Firstname         string `json:"first_name"`
	Lastname          string `json:"last_name"`
	Mobile            string `json:"mobile"`
	Email             string `json:"email"`
}
type AccountResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	profile := internal.Profile{
		Profile_Picture:    req.ProfilePicure,
		Profile_Background: req.ProfileBackground,
		First_Name:         req.Firstname,
		Last_Name:          req.Lastname,
		Mobile:             req.Mobile,
		Email:              req.Email,
	}
	_, err := rb.svc.CreateAccount(r.Context(), internal.Account{
		UserName: req.Username,
		Profile:  profile,
	}, req.Password)
	if err != nil {
		renderErrorResponse(r.Context(), w, "create failed", err)
		return
	}
	renderResponse(w,
		&AccountResponse{
			Message: "Created Successfully",
		}, http.StatusCreated)
}

type ReadAccountResponse struct {
	Account Account `json:"account"`
}

func (rb *RBACHandler) account(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	account, err := rb.svc.Account(r.Context(), username)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	profile := Profile{
		Id:                account.Profile.Id,
		ProfileBackground: account.Profile.Profile_Background,
		ProfilePicture:    account.Profile.Profile_Picture,
		FirstName:         account.Profile.First_Name,
		LastName:          account.Profile.Last_Name,
		Mobile:            account.Profile.Mobile,
		Email:             account.Profile.Email,
		CreatedAt:         account.CreatedAt,
	}
	renderResponse(w, &ReadAccountResponse{
		Account: Account{
			Id:        account.Id,
			Username:  account.UserName,
			Profile:   profile,
			CreatedAt: account.CreatedAt,
		},
	}, http.StatusOK)
}

type ListAccountRequest struct {
	From int `json:"from"`
	Size int `json:"size"`
}
type ListAccountResponse struct {
	Accounts []Account `json:"accounts"`
	Total    int64     `json:"total"`
}

func (rb *RBACHandler) listaccount(w http.ResponseWriter, r *http.Request) {
	var req ListAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	la, err := rb.svc.ListAccount(r.Context(), internal.ListArgs{
		From: &req.From,
		Size: &req.Size,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	accounts := []Account{}
	for _, value := range la.Accounts {
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
		accounts = append(accounts, acc)
	}
	renderResponse(w, &ListAccountResponse{
		Accounts: accounts,
		Total:    la.Total,
	}, http.StatusOK)
}

type UpdateProfileRequest struct {
	Id                string `json:"id"`
	ProfilePicture    string `json:"profile_picture"`
	ProfileBackground string `json:"profile_background"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Mobile            string `json:"mobile"`
	Email             string `json:"email"`
}

func (rb *RBACHandler) updateProfile(w http.ResponseWriter, r *http.Request) {
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.UpdateProfile(r.Context(), internal.Profile{
		Id:                 req.Id,
		Profile_Picture:    req.ProfilePicture,
		Profile_Background: req.ProfileBackground,
		First_Name:         req.FirstName,
		Last_Name:          req.LastName,
		Mobile:             req.Mobile,
		Email:              req.Email,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating profile", err)
		return
	}
	renderResponse(w,
		&AccountResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}

type DeleteAccountResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) deleteAccount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	_, err := rb.svc.Account(r.Context(), username)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	err = rb.svc.DeleteAccount(r.Context(), username)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error deleting the account", err)
		return
	}

	renderResponse(w, &DeleteAccountResponse{
		Message: "Deleted Successfully..",
	}, http.StatusOK)
}

type AccountRoleByAccount struct {
	Account Account `json:"account"`
	Roles   []Role  `json:"roles"`
}

func (rb *RBACHandler) getAccountRoleByAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["username"]
	la, err := rb.svc.AccountRoleByAccount(r.Context(), id)
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	prof := Profile{
		Id:                la.Account.Profile.Id,
		ProfilePicture:    la.Account.Profile.Profile_Picture,
		ProfileBackground: la.Account.Profile.Profile_Background,
		FirstName:         la.Account.Profile.First_Name,
		LastName:          la.Account.Profile.Last_Name,
		Mobile:            la.Account.Profile.Mobile,
		Email:             la.Account.Profile.Email,
		CreatedAt:         la.Account.CreatedAt,
	}
	acc := Account{
		Id:        la.Account.Id,
		Username:  la.Account.UserName,
		Profile:   prof,
		CreatedAt: la.Account.CreatedAt,
	}
	roles := []Role{}
	for _, value := range la.Roles {
		//get role
		rl, err := rb.svc.Role(r.Context(), value.Id)
		if err != nil {
			renderErrorResponse(r.Context(), w, "error getting role", err)
			return
		}
		roles = append(roles, Role{
			Id:        rl.Id,
			Role:      rl.Role,
			CreatedAt: rl.CreatedAt,
		})
	}
	renderResponse(w, &AccountRoleByAccount{
		Account: acc,
		Roles:   roles,
	}, http.StatusOK)
}
