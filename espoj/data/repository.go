package data

import (
	"sort"
	"time"
)

type Commit struct {
	Sha       string
	ParentSha *[]string
	Time      *time.Time
	ShowUrl   string
	Message   *string
}
type Issue struct {
}
type PullRequest struct {
}
type Repository struct {
	Url, Name, Description string
	Commits                []Commit
}

func (repo *Repository) SortCommits() {
	sort.Slice(repo.Commits, func(i, j int) bool {
		return repo.Commits[i].Time.Before(*repo.Commits[j].Time)
	})
}
