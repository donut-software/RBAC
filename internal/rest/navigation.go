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
	TaskId    string    `json:"taskId"`
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
			TaskId:    navigation.Task_id,
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

type ListNavigationRequest struct {
	From int `json:"from"`
	Size int `json:"size"`
}
type ListNavigationResponse struct {
	Navigation []Navigation `json:"navigations"`
	Total      int64        `json:"total"`
}

func (rb *RBACHandler) listNavigation(w http.ResponseWriter, r *http.Request) {
	var req ListNavigationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	la, err := rb.svc.ListNavigation(r.Context(), internal.ListArgs{
		From: &req.From,
		Size: &req.Size,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	helpText := []Navigation{}
	for _, value := range la.Navigation {
		helpText = append(helpText, Navigation{
			Id:        value.Id,
			Name:      value.Name,
			TaskId:    value.Task_id,
			CreatedAt: value.CreatedAt,
		})
	}
	renderResponse(w, &ListNavigationResponse{
		Navigation: helpText,
		Total:      la.Total,
	}, http.StatusOK)
}
