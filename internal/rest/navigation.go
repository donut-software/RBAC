package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
	"time"

	"github.com/gorilla/mux"
)

type Navigation struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateNavigationRequest struct {
	Name   string `json:"name"`
	TaskId string `json:"taskId"`
}
type NavigationResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) createNavigation(w http.ResponseWriter, r *http.Request) {
	var req CreateNavigationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.CreateNavigation(r.Context(), internal.Navigation{
		Name:    req.Name,
		Task_id: req.TaskId,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "create navigation failed", err)
		return
	}
	renderResponse(w,
		&NavigationResponse{
			Message: "Created Successfully",
		}, http.StatusCreated)
}

type GetNavigationResponse struct {
	Navigation Navigation `json:"navigation"`
}

func (rb *RBACHandler) navigation(w http.ResponseWriter, r *http.Request) {
	navigationId := mux.Vars(r)["navigationId"]
	navigation, err := rb.svc.Navigation(r.Context(), navigationId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	renderResponse(w, &GetNavigationResponse{
		Navigation: Navigation{
			Id:        navigation.Id,
			Name:      navigation.Name,
			CreatedAt: navigation.CreatedAt,
		},
	}, http.StatusOK)
}

type UpdateNavigationRequest struct {
	NavigationId string `json:"navigationId"`
	Name         string `json:"name"`
	TaskId       string `json:"taskId"`
}

func (rb *RBACHandler) updateNavigation(w http.ResponseWriter, r *http.Request) {
	var req UpdateNavigationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.UpdateNavigation(r.Context(), internal.Navigation{
		Id:      req.NavigationId,
		Name:    req.Name,
		Task_id: req.TaskId,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating navigation", err)
		return
	}
	renderResponse(w,
		&NavigationResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}
