package github

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hasyimibhar/github-search/common"
)

func TestSearchGateway_Repositories(t *testing.T) {
	client := NewClient(&ClientConfig{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}, common.NewStandardLogger(os.Stderr, "trace"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := SearchRepositoriesRequest{}.
		Language("go")

	resp, err := client.Search().Repositories(ctx, req, PaginationOptions{Page: 1, PerPage: 1})
	if err != nil {
		t.Fatal(err)
	}

	if resp.TotalCount == 0 {
		t.Fatal("expecting at least 1 item")
	}

	if len(resp.Items) != 1 {
		t.Fatal("expecting 1 item per page")
	}

	repo := resp.Items[0]
	if repo.FullName != "golang/go" {
		t.Fatal("search result should return 'golang/go' first")
	}
}

func TestSearchGateway_Repositories_MultipleLanguages(t *testing.T) {
	client := NewClient(&ClientConfig{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}, common.NewStandardLogger(os.Stderr, "trace"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := SearchRepositoriesRequest{}.
		Language("go").Language("ruby")

	resp, err := client.Search().Repositories(ctx, req, PaginationOptions{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatal(err)
	}

	if resp.TotalCount == 0 {
		t.Fatal("expecting at least 1 item")
	}
}

func TestSearchGateway_Repositories_CPlusPlus(t *testing.T) {
	client := NewClient(&ClientConfig{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}, common.NewStandardLogger(os.Stderr, "trace"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := SearchRepositoriesRequest{}.
		Language("C++")

	resp, err := client.Search().Repositories(ctx, req, PaginationOptions{Page: 1, PerPage: 1})
	if err != nil {
		t.Fatal(err)
	}

	if resp.TotalCount == 0 {
		t.Fatal("expecting at least 1 item")
	}

	if len(resp.Items) != 1 {
		t.Fatal("expecting 1 item per page")
	}

	repo := resp.Items[0]
	if repo.Language != "C++" {
		t.Fatal("search result should return C++, not C")
	}
}

func TestSearchGateway_RepositoriesPagination(t *testing.T) {
	client := NewClient(&ClientConfig{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}, common.NewStandardLogger(os.Stderr, "trace"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tests := []struct {
		PerPage int
	}{
		{5},
		{10},
	}

	for _, tt := range tests {
		req := SearchRepositoriesRequest{}.
			Language("Go")

		resp, err := client.Search().Repositories(ctx, req, PaginationOptions{Page: 1, PerPage: tt.PerPage})
		if err != nil {
			t.Fatal(err)
		}

		if len(resp.Items) != tt.PerPage {
			t.Fatalf("expecting %d items per page", tt.PerPage)
		}
	}
}
