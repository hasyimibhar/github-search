package github

import (
	"net/http"
	"strconv"
)

// RequestOptions modifies the request.
type RequestOptions interface {
	Apply(req *http.Request)
}

// PaginationOptions adds pagination to the request.
type PaginationOptions struct {
	Page    int
	PerPage int
}

func (o PaginationOptions) Apply(req *http.Request) {
	q := req.URL.Query()
	if o.Page > 0 {
		q.Set("page", strconv.Itoa(o.Page))
	}
	if o.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(o.PerPage))
	}

	req.URL.RawQuery = q.Encode()
}

// AuthenticationOptions adds authentication to the request.
type AuthenticationOptions struct {
	ClientID     string
	ClientSecret string
}

func (o AuthenticationOptions) Apply(req *http.Request) {
	q := req.URL.Query()
	q.Set("client_id", o.ClientID)
	q.Set("client_secret", o.ClientSecret)

	req.URL.RawQuery = q.Encode()
}
