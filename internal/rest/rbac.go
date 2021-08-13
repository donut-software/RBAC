package rest

import (
	"context"
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
	r.HandleFunc("/accounts/{username}", rb.account).Methods(http.MethodGet)
}
