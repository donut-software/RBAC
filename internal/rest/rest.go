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
