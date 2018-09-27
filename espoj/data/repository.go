package data

import (
	"os"
	"os/exec"
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
	CloneUrl               string
	Commits                []Commit
}

func (repo *Repository) SortCommits() {
	sort.Slice(repo.Commits, func(i, j int) bool {
		return repo.Commits[i].Time.Before(*repo.Commits[j].Time)
	})
}
func (repo *Repository) Save() {
	repoPath := "repos" + string(os.PathSeparator) + repo.Name
	_ = os.MkdirAll(repoPath, os.ModePerm)
	cmd := exec.Command("git", "clone", "--mirror", repo.CloneUrl, repoPath)
	_ = cmd.Start()
	cmd.Wait()
}
