package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	apihttp "github.com/hasyimibhar/github-search/api/http"
	"github.com/hasyimibhar/github-search/report"
)

type ReportController struct {
	Database *report.Database
}

func (c *ReportController) GetReport(w http.ResponseWriter, r *http.Request) {
	from := timeParam(r, "from")
	to := timeParam(r, "to")

	if from.Equal(to) || from.After(to) {
		render.Render(w, r, apihttp.ErrInvalidRequest(fmt.Errorf("invalid time range")))
		return
	}

	report, err := c.Database.GenerateReport(from, to)
	if err != nil {
		render.Render(w, r, apihttp.ErrInternalServerError(
			fmt.Errorf("failed to generate report: %s", err)))
		return
	}

	render.Render(w, r, SearchReportView{
		From:               report.From.Format("2006-01-02 15:04:05"),
		To:                 report.To.Format("2006-01-02 15:04:05"),
		TotalSearches:      report.TotalSearches,
		FailedSearches:     report.FailedSearches,
		SearchesByTopic:    report.SearchesByTopic,
		SearchesByLanguage: report.SearchesByLanguage,
	})
}

func timeParam(r *http.Request, param string) time.Time {
	s := r.URL.Query().Get(param)
	i, err := strconv.Atoi(s)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(int64(i), 0).UTC()
}

type SearchReportView struct {
	From               string         `json:"from"`
	To                 string         `json:"to"`
	TotalSearches      int            `json:"total_searches"`
	FailedSearches     int            `json:"failed_searches"`
	SearchesByTopic    map[string]int `json:"searches_by_topic"`
	SearchesByLanguage map[string]int `json:"searches_by_language"`
}

func (v SearchReportView) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
