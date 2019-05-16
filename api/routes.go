package api

import (
	"github.com/go-chi/chi"
	"github.com/hasyimibhar/github-search/api/admin"
	"github.com/hasyimibhar/github-search/api/search"
	"github.com/hasyimibhar/github-search/common"
	"github.com/hasyimibhar/github-search/github"
	"github.com/hasyimibhar/github-search/report"
	"upper.io/db.v3/lib/sqlbuilder"
)

func Router(githubClient *github.Client, database sqlbuilder.Database, log common.Logger) chi.Router {
	r := chi.NewRouter()

	reportDatabase := &report.Database{
		Database: database,
		Log:      log,
	}

	r.Mount("/search", search.Router(&search.SearchController{
		GithubClient: githubClient,
		SearchLogger: &report.Logger{Database: database},
		Log:          log,
	}))

	r.Mount("/admin", admin.Router(&admin.ReportController{
		Database: reportDatabase,
	}, &admin.SearchController{
		Database: reportDatabase,
	}))

	return r
}
