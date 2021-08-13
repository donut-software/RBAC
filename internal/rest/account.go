package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
	"time"

	"github.com/gorilla/mux"
)

type Profile struct {
	ProfilePicture    string    `json:"profile_picture"`
	ProfileBackground string    `json:"profile_background"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Mobile            string    `json:"mobile"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
}

type Account struct {
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
			Username:  account.UserName,
			Profile:   profile,
			CreatedAt: account.CreatedAt,
		},
	}, http.StatusOK)
}
