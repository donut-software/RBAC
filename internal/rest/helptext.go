package rest

import (
	"encoding/json"
	"net/http"
	"rbac/internal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type HelpText struct {
	Id        string    `json:"id"`
	HelpText  string    `json:"helpText"`
	TaskId    string    `json:"taskId"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateHelpTextRequest struct {
	TaskId   string `json:"taskId"`
	HelpText string `json:"helptext"`
}
type HelpTextResponse struct {
	Message string `json:"message"`
}

func (rb *RBACHandler) createHelpText(w http.ResponseWriter, r *http.Request) {
	var req CreateHelpTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.CreateHelpText(r.Context(), internal.HelpText{
		HelpText: req.HelpText,
		Task_id:  req.TaskId,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "create helpText failed", err)
		return
	}
	renderResponse(w,
		&HelpTextResponse{
			Message: "Created Successfully",
		}, http.StatusCreated)
}

type GetHelpTextResponse struct {
	HelpText HelpText `json:"helpText"`
}

func (rb *RBACHandler) helpText(w http.ResponseWriter, r *http.Request) {
	helpTextId := mux.Vars(r)["helpTextId"]
	helpText, err := rb.svc.HelpText(r.Context(), helpTextId)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error getting the account", err)
		return
	}
	renderResponse(w, &GetHelpTextResponse{
		HelpText: HelpText{
			Id:        helpText.Id,
			HelpText:  helpText.HelpText,
			TaskId:    helpText.Task_id,
			CreatedAt: helpText.CreatedAt,
		},
	}, http.StatusOK)
}

type UpdateHelpTextRequest struct {
	HelpTextId string `json:"helpTextId"`
	HelpText   string `json:"helpText"`
	TaskId     string `json:"taskId"`
}

func (rb *RBACHandler) updateHelpText(w http.ResponseWriter, r *http.Request) {
	var req UpdateHelpTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	err := rb.svc.UpdateHelpText(r.Context(), internal.HelpText{
		Id:       req.HelpTextId,
		HelpText: req.HelpText,
		Task_id:  req.TaskId,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "error updating helpText", err)
		return
	}
	renderResponse(w,
		&HelpTextResponse{
			Message: "Updated Successfully",
		}, http.StatusCreated)
}

type ListHelpTextResponse struct {
	HelpText []HelpText `json:"helptexts"`
	Total    int64      `json:"total"`
}

func (rb *RBACHandler) listHelpText(w http.ResponseWriter, r *http.Request) {

	var from int
	var size int
	v := r.URL.Query()
	from, err := strconv.Atoi(v.Get("from"))
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid param from", err)
		return
	}
	size, err = strconv.Atoi(v.Get("size"))
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid param size", err)
		return
	}
	la, err := rb.svc.ListHelpText(r.Context(), internal.ListArgs{
		From: &from,
		Size: &size,
	})
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	helpText := []HelpText{}
	for _, value := range la.HelpText {
		helpText = append(helpText, HelpText{
			Id:        value.Id,
			HelpText:  value.HelpText,
			TaskId:    value.Task_id,
			CreatedAt: value.CreatedAt,
		})
	}
	renderResponse(w, &ListHelpTextResponse{
		HelpText: helpText,
		Total:    la.Total,
	}, http.StatusOK)
}

func (rb *RBACHandler) deleteHelpText(w http.ResponseWriter, r *http.Request) {
	htID := mux.Vars(r)["helpTextId"]
	err := rb.svc.DeleteHelpText(r.Context(), htID)
	if err != nil {
		renderErrorResponse(r.Context(), w, "error deleting helptext", err)
		return
	}
	renderResponse(w,
		&HelpTextResponse{
			Message: "Deleted Successfully",
		}, http.StatusOK)
}
