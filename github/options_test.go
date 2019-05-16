package github

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPaginationOptions(t *testing.T) {
	tests := []struct {
		Page    int
		PerPage int
	}{
		{1, 1},
		{2, 5},
		{3, 10},
		{4, 20},
		{5, 50},
		{6, 100},
		{0, 1},
		{1, 0},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			opt := PaginationOptions{
				Page:    tt.Page,
				PerPage: tt.PerPage,
			}

			req := request(t)
			opt.Apply(req)

			var expectedURL string
			if tt.Page > 0 && tt.PerPage > 0 {
				expectedURL = fmt.Sprintf("%s/search/repositories?page=%d&per_page=%d", BaseURL, tt.Page, tt.PerPage)
			} else if tt.Page > 0 {
				expectedURL = fmt.Sprintf("%s/search/repositories?page=%d", BaseURL, tt.Page)
			} else if tt.PerPage > 0 {
				expectedURL = fmt.Sprintf("%s/search/repositories?per_page=%d", BaseURL, tt.PerPage)
			}

			if req.URL.String() != expectedURL {
				t.Fatal("url does not match")
			}
		})
	}
}

func TestAuthenticationOptions(t *testing.T) {
	tests := []struct {
		ClientID     string
		ClientSecret string
	}{
		{"foor", "bar"},
		{"lorem", "ipsum"},
		{"", ""},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			opt := AuthenticationOptions{
				ClientID:     tt.ClientID,
				ClientSecret: tt.ClientSecret,
			}

			req := request(t)
			opt.Apply(req)

			expectedURL := fmt.Sprintf("%s/search/repositories?client_id=%s&client_secret=%s", BaseURL, tt.ClientID, tt.ClientSecret)
			if req.URL.String() != expectedURL {
				t.Fatal("url does not match")
			}
		})
	}
}

func request(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/search/repositories", BaseURL), nil)
	if err != nil {
		t.Fatal(err)
	}

	return req
}
