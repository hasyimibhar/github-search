package report

import (
	"context"
	"strings"
	"time"

	"upper.io/db.v3/lib/sqlbuilder"
)

type Logger struct {
	Database sqlbuilder.Database
}

type SearchEntry struct {
	Topics          []string
	Languages       []string
	ResponseStatus  int
	ResponseContent string
	CreatedAt       time.Time
}

// LogSearch logs a search.
func (l *Logger) LogSearch(ctx context.Context, e SearchEntry) error {
	_, err := l.Database.WithContext(ctx).Collection("searches").Insert(makeSearchDto(e))
	return err
}

type searchDto struct {
	Topics          string    `db:"topics"`
	Languages       string    `db:"languages"`
	ResponseStatus  int       `db:"response_status"`
	ResponseContent string    `db:"response_content"`
	CreatedAt       time.Time `db:"created_at"`
}

func makeSearchDto(e SearchEntry) searchDto {
	return searchDto{
		Topics:          strings.ToLower(strings.Join(e.Topics, ",")),
		Languages:       strings.ToLower(strings.Join(e.Languages, ",")),
		ResponseStatus:  e.ResponseStatus,
		ResponseContent: e.ResponseContent,
		CreatedAt:       e.CreatedAt,
	}
}

func (d searchDto) toSearchEntry() SearchEntry {
	topics := []string{}
	if d.Topics != "" {
		topics = strings.Split(d.Topics, ",")
	}

	languages := []string{}
	if d.Languages != "" {
		languages = strings.Split(d.Languages, ",")
	}

	return SearchEntry{
		Topics:          topics,
		Languages:       languages,
		ResponseStatus:  d.ResponseStatus,
		ResponseContent: d.ResponseContent,
		CreatedAt:       d.CreatedAt,
	}
}
