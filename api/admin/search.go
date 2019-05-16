package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	apihttp "github.com/hasyimibhar/github-search/api/http"
	"github.com/hasyimibhar/github-search/report"
)

const (
	DefaultLatestSearches = 10
)

type SearchController struct {
	Database *report.Database
}

// LatestSearches returns the latest N searches.
func (c *SearchController) LatestSearches(w http.ResponseWriter, r *http.Request) {
	n := intParam(r, "n")
	if n == 0 {
		n = DefaultLatestSearches
	}

	searches, err := c.Database.LatestSearches(n)
	if err != nil {
		render.Render(w, r, apihttp.ErrInternalServerError(
			fmt.Errorf("failed to query searches: %s", err)))
		return
	}

	views := []render.Renderer{}

	for _, s := range searches {
		views = append(views, SearchEntryView{
			Topics:          s.Topics,
			Languages:       s.Languages,
			ResponseStatus:  s.ResponseStatus,
			ResponseContent: s.ResponseContent,
			CreatedAt:       s.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	render.RenderList(w, r, views)
}

func intParam(r *http.Request, param string) int {
	s := r.URL.Query().Get(param)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}

type SearchEntryView struct {
	Topics          []string `json:"topics"`
	Languages       []string `json:"languages"`
	ResponseStatus  int      `json:"response_status"`
	ResponseContent string   `json:"response_content"`
	CreatedAt       string   `json:"created_at"`
}

func (v SearchEntryView) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
