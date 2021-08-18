package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rbac/internal"
)

type errorResponse struct {
	Error string `json:"error"`
}

func renderErrorResponse(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	resp := errorResponse{Error: msg}
	status := http.StatusInternalServerError

	var ierr *internal.Error
	if !errors.As(err, &ierr) {
		resp.Error = "internal error"
	} else {
		switch ierr.Code() {
		case internal.ErrorCodeNotFound:
			status = http.StatusNotFound
		case internal.ErrorCodeInvalidArgument:
			status = http.StatusBadRequest
		}
	}
	if err != nil {
		// _, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "rest.renderErrorResponse")
		// defer span.End()
		// span.RecordError(err)
		fmt.Println(err)
	}
	renderResponse(w, resp, status)
}

func renderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		// XXX Do something with the error ;)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		// XXX Do something with the error ;)
		fmt.Println(err)
	}
}

// func convertInternalNavigationList(list []internal.Navigation) []Navigation {
// 	navlist := []Navigation{}
// 	for _, value := range list {
// 		n := Navigation{
// 			Id:        value.Id,
// 			Name:      value.Name,
// 			TaskId:    value.Task_id,
// 			CreatedAt: value.CreatedAt,
// 		}
// 		navlist = append(navlist, n)
// 	}
// 	return navlist
// }
// func convertInternalMenuList(list []internal.Menu) []Menu {
// 	menulist := []Menu{}
// 	for _, value := range list {
// 		m := Menu{
// 			Id:        value.Id,
// 			Name:      value.Name,
// 			TaskId:    value.Task_id,
// 			CreatedAt: value.CreatedAt,
// 		}
// 		menulist = append(menulist, m)
// 	}
// 	return menulist
// }

func addCookie(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
