package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"

	ErrInternal = "Internal server error"
	ErrNotFound = "not found"
	ErrBadReq   = "Bad request"
	ErrWrongID  = "Wrong ID"
)

func Ok() Response {
	return Response{
		Status: StatusOK,
	}
}

func SendErrorInternal(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, Response{
		Status: StatusError,
		Error:  ErrInternal,
	})
}

func SendErrorBadReq(r *http.Request, w http.ResponseWriter, msg string) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, Response{
		Status:  StatusError,
		Error:   ErrBadReq,
		Message: msg,
	})
}

func SendErrorNotFound(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, Response{
		Status: StatusError,
		Error:  ErrNotFound,
	})
}
