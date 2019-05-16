package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hasyimibhar/github-search/common"
)

type SearchGateway struct {
	HTTPClient *http.Client

	config *ClientConfig
	logger common.Logger
}

// SearchRepositoriesRequest contains the search criteria for searching
// repositories.
type SearchRepositoriesRequest struct {
	Topics    []string
	Languages []string
}

func (r SearchRepositoriesRequest) Topic(topic string) SearchRepositoriesRequest {
	r.Topics = append(r.Topics, topic)
	return r
}

func (r SearchRepositoriesRequest) Language(language string) SearchRepositoriesRequest {
	r.Languages = append(r.Languages, language)
	return r
}

func (r SearchRepositoriesRequest) Encode() string {
	params := []string{}
	for _, t := range r.Topics {
		params = append(params, "topic:"+t)
	}
	for _, l := range r.Languages {
		params = append(params, "language:"+l)
	}

	return strings.Join(params, " ")
}

// Owner is the owner of a Repository.
type Owner struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"url"`
}

// Repository is a Github repository.
type Repository struct {
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	HTMLURL         string    `json:"html_url"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	ForksCount      int       `json:"forks_count"`
	Language        string    `json:"language"`
	Score           float64   `json:"score"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SearchRepositoriesResponse is the response of calling client.Search().Repositories(...).
type SearchRepositoriesResponse struct {
	Items      []Repository `json:"items"`
	TotalCount int          `json:"total_count"`
}

// Repositories returns the repositories which match the search criteria.
func (g SearchGateway) Repositories(
	ctx context.Context,
	req SearchRepositoriesRequest,
	opts ...RequestOptions) (*SearchRepositoriesResponse, error) {

	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/search/repositories", BaseURL), nil)
	if err != nil {
		return nil, err
	}

	q := httpReq.URL.Query()
	q.Set("q", req.Encode())
	httpReq.URL.RawQuery = q.Encode()

	if g.config != nil {
		opts = append(opts, AuthenticationOptions{
			ClientID:     g.config.ClientID,
			ClientSecret: g.config.ClientSecret,
		})
	}

	for _, o := range opts {
		o.Apply(httpReq)
	}

	httpReq = httpReq.WithContext(ctx)
	httpReq.Header.Set("Accept", "application/vnd.github.mercy-preview+json")

	g.logger.Tracef("sending HTTP GET request to %s", httpReq.URL.String())

	httpResp, err := g.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		var resp Error
		resp.StatusCode = httpResp.StatusCode

		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %s", err)
		}

		return nil, resp
	}

	defer httpResp.Body.Close()

	var resp SearchRepositoriesResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %s", err)
	}

	return &resp, nil
}
