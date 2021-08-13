package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
	"time"

	"github.com/gorilla/mux"
)

type Menu struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateMenuRequest struct {
	Name   string `json:"name"`
	TaskId string `json:"taskId"`
}
type MenuResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) createMenu(w http.ResponseWriter, r *http.Request) {
	var req CreateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.CreateMenu(r.Context(), internal.Menu{
		Name:    req.Name,
		Task_id: req.TaskId,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "create menu failed", err)
		return
	}
	renderResponse(w,
		&MenuResponse{
			Message: "Created Successfully",
		}, http.StatusCreated)
}

type GetMenuResponse struct {
	Menu Menu `json:"menu"`
}

func (rb *RBACHandler) menu(w http.ResponseWriter, r *http.Request) {
	menuId := mux.Vars(r)["menuId"]
	menu, err := rb.svc.Menu(r.Context(), menuId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	renderResponse(w, &GetMenuResponse{
		Menu: Menu{
			Id:        menu.Id,
			Name:      menu.Name,
			CreatedAt: menu.CreatedAt,
		},
	}, http.StatusOK)
}

type UpdateMenuRequest struct {
	MenuId string `json:"menuId"`
	Name   string `json:"name"`
	TaskId string `json:"taskId"`
}

func (rb *RBACHandler) updateMenu(w http.ResponseWriter, r *http.Request) {
	var req UpdateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.UpdateMenu(r.Context(), internal.Menu{
		Id:      req.MenuId,
		Name:    req.Name,
		Task_id: req.TaskId,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating menu", err)
		return
	}
	renderResponse(w,
		&MenuResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}
