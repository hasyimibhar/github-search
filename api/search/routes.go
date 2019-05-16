package search

import (
	"github.com/go-chi/chi"
)

func Router(search *SearchController) chi.Router {
	r := chi.NewRouter()

	r.Get("/", search.Search)

	return r
}
