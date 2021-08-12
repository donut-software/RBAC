package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"rbac/internal"

	"github.com/gorilla/mux"
)

//go:generate counterfeiter -o resttesting/rbac_service.gen.go . RBACService
type RBACService interface {
	CreateAccount(ctx context.Context, account internal.Account, password string) error
	Account(ctx context.Context, username string) (internal.Account, error)
	UpdateProfile(ctx context.Context, profile internal.Profile) error
	ChangePassword(ctx context.Context, username string, password string) error
	DeleteAccount(ctx context.Context, username string) error
}

type RBACHandler struct {
	svc RBACService
}

func NewRBACHandler(svc RBACService) *RBACHandler {
	return &RBACHandler{
		svc: svc,
	}
}

func (rb *RBACHandler) Register(r *mux.Router) {
	r.HandleFunc("/register", rb.register).Methods(http.MethodPost)
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
	err := rb.svc.CreateAccount(r.Context(), internal.Account{
		UserName: req.Username,
		Profile:  profile,
	}, req.Password)
	if err != nil {
		renderErrorResponse(r.Context(), w, "create failed", err)
		return
	}
	renderResponse(w,
		&AccountResponse{
			Message: "Created Succesfully",
		}, http.StatusCreated)
}
