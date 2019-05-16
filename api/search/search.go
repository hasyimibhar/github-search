package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/render"
	apihttp "github.com/hasyimibhar/github-search/api/http"
	"github.com/hasyimibhar/github-search/common"
	"github.com/hasyimibhar/github-search/github"
	"github.com/hasyimibhar/github-search/report"
)

type SearchController struct {
	GithubClient *github.Client
	SearchLogger *report.Logger
	Log          common.Logger
}

func (c *SearchController) Search(w http.ResponseWriter, r *http.Request) {
	topics := stringsParam(r, "topics")
	languages := stringsParam(r, "languages")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := c.GithubClient.Search().Repositories(ctx, github.SearchRepositoriesRequest{
		Topics:    topics,
		Languages: languages,
	}, github.PaginationOptions{
		Page:    intParam(r, "page"),
		PerPage: intParam(r, "per_page"),
	})

	go c.logSearch(topics, languages, resp, err)

	if err != nil {
		githubErr, ok := err.(github.Error)
		if ok {
			w.WriteHeader(githubErr.StatusCode)
			render.Render(w, r, ErrorView{Message: githubErr.Message})
			return
		} else {
			render.Render(w, r, apihttp.ErrInternalServerError(err))
			return
		}
	}

	render.Render(w, r, SearchResultView{
		TotalCount: resp.TotalCount,
		Items:      resp.Items,
	})
}

func (c *SearchController) logSearch(
	topics []string, languages []string,
	resp *github.SearchRepositoriesResponse,
	err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	entry := report.SearchEntry{
		Topics:    topics,
		Languages: languages,
		CreatedAt: time.Now(),
	}

	if err != nil {
		githubErr, ok := err.(github.Error)
		if ok {
			entry.ResponseStatus = githubErr.StatusCode
			js, _ := json.Marshal(githubErr)
			entry.ResponseContent = string(js)
		}
	} else {
		entry.ResponseStatus = http.StatusOK
		js, _ := json.Marshal(resp)
		entry.ResponseContent = string(js)
	}

	if err := c.SearchLogger.LogSearch(ctx, entry); err != nil {
		c.Log.Error(fmt.Errorf("failed to log search: %s", err))
	}
}

func stringsParam(r *http.Request, param string) []string {
	s := r.URL.Query().Get(param)

	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

func intParam(r *http.Request, param string) int {
	s := r.URL.Query().Get(param)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}

type ErrorView struct {
	Message string `json:"message"`
}

func (v ErrorView) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type SearchResultView struct {
	TotalCount int                 `json:"total_count"`
	Items      []github.Repository `json:"items"`
}

func (v SearchResultView) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
