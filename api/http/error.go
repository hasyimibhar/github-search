package http

import (
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrNotFound     = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found"}
	ErrUnauthorized = &ErrResponse{HTTPStatusCode: 401, StatusText: "Unauthorized"}
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

func ErrUnprocessableEntity(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Unprocessable entity",
		ErrorText:      err.Error(),
	}
}

func ErrForbidden(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 403,
		StatusText:     "Forbidden",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error",
		ErrorText:      err.Error(),
	}
}
