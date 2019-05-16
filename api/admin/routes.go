package admin

import (
	"github.com/go-chi/chi"
)

func Router(report *ReportController, search *SearchController) chi.Router {
	r := chi.NewRouter()

	r.Mount("/reports", ReportRouter(report))
	r.Mount("/searches", SearchRouter(search))

	return r
}

func ReportRouter(report *ReportController) chi.Router {
	r := chi.NewRouter()

	r.Get("/", report.GetReport)

	return r
}

func SearchRouter(search *SearchController) chi.Router {
	r := chi.NewRouter()

	r.Get("/", search.LatestSearches)

	return r
}
