package report

import (
	"net/http"
	"time"

	"github.com/hasyimibhar/github-search/common"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type SearchReport struct {
	From time.Time
	To   time.Time

	TotalSearches      int
	FailedSearches     int
	SearchesByTopic    map[string]int
	SearchesByLanguage map[string]int
}

// Database contains all the searches that has been done.
type Database struct {
	Database sqlbuilder.Database
	Log      common.Logger
}

// LatestSearches returns the latest N searches.
func (d *Database) LatestSearches(n int) (searches []SearchEntry, err error) {
	q := d.Database.Collection("searches").Find().OrderBy("-created_at").Limit(n)

	var dtos []searchDto
	err = q.All(&dtos)

	searches = []SearchEntry{}
	for _, d := range dtos {
		searches = append(searches, d.toSearchEntry())
	}

	return
}

// GenerateReport generates a search report based on the time range specified.
func (d *Database) GenerateReport(from time.Time, to time.Time) (*SearchReport, error) {
	searches, err := d.searches(from, to)
	if err != nil {
		return nil, err
	}

	report := SearchReport{From: from, To: to}
	report.TotalSearches = len(searches)
	report.SearchesByTopic = map[string]int{}
	report.SearchesByLanguage = map[string]int{}

	for _, s := range searches {
		// Ignore failed requests for the report
		if s.ResponseStatus != http.StatusOK {
			report.FailedSearches++
			continue
		}

		for _, t := range s.Topics {
			report.SearchesByTopic[t] = report.SearchesByTopic[t] + 1
		}
		for _, l := range s.Languages {
			report.SearchesByLanguage[l] = report.SearchesByLanguage[l] + 1
		}
	}

	return &report, nil
}

// Searches returns the searches done between the specified time range.
func (d *Database) searches(from time.Time, to time.Time) (searches []SearchEntry, err error) {
	q := d.Database.Collection("searches").Find().Where(db.Cond{
		"created_at >=": from.Format("2006-01-02 15:04:05"),
		"created_at <=": to.Format("2006-01-02 15:04:05"),
	})

	var dtos []searchDto
	err = q.All(&dtos)

	searches = []SearchEntry{}
	for _, d := range dtos {
		searches = append(searches, d.toSearchEntry())
	}

	return
}
