package github

import (
	"context"
	"github.com/google/go-github/v18/github"
	"github.com/samitc/espoj/espoj/data"
	"strings"
)

type Importer struct {
	Url        string
	client     *github.Client
	repository *github.Repository
	commits    []*github.RepositoryCommit
}

func extractNameAndOwner(url string) (name, owner string) {
	nameIndex := strings.LastIndex(url, "/")
	name = url[nameIndex+1:]
	ownerIndex := strings.LastIndex(url[0:nameIndex], "/")
	owner = url[ownerIndex+1 : nameIndex]
	return
}
func loadCommits(client *github.Client, owner, name string) []*github.RepositoryCommit {
	const perPage = 2000
	commits, resp, _ := client.Repositories.ListCommits(context.Background(), owner, name, &github.CommitsListOptions{ListOptions: github.ListOptions{PerPage: perPage}})
	var pCommits []*github.RepositoryCommit
	for resp.NextPage != 0 {
		pCommits, resp, _ = client.Repositories.ListCommits(context.Background(), owner, name, &github.CommitsListOptions{ListOptions: github.ListOptions{PerPage: perPage, Page: resp.NextPage}})
		commits = append(commits, pCommits...)
	}
	return commits
}
func (ghi *Importer) initClient() {
	if ghi.client == nil {
		ghi.client = github.NewClient(nil)
		name, owner := extractNameAndOwner(ghi.Url)
		ghi.repository, _, _ = ghi.client.Repositories.Get(context.Background(), owner, name)
		ghi.commits = loadCommits(ghi.client, owner, name)
	}
}
func (ghi *Importer) GetUrl() string {
	ghi.initClient()
	return ghi.repository.GetHTMLURL()
}
func (ghi *Importer) GetName() string {
	ghi.initClient()
	return ghi.repository.GetName()
}
func (ghi *Importer) GetDescription() string {
	ghi.initClient()
	return ghi.repository.GetDescription()
}
func (ghi *Importer) GetCommits() []data.Commit {
	ghi.initClient()
	commits := make([]data.Commit, len(ghi.commits))
	for i, commit := range ghi.commits {
		com := commit.Commit
		parentSha := make([]string, len(commit.Parents))
		for i, par := range commit.Parents {
			parentSha[i] = *par.SHA
		}
		commits[i] = data.Commit{Sha: *commit.SHA, ShowUrl: *commit.HTMLURL, ParentSha: &parentSha, Time: com.Author.Date,
			Message: com.Message}
	}
	return commits
}
func (ghi *Importer) GetCloneUrl() string {
	ghi.initClient()
	return ghi.repository.GetCloneURL()
}
func (ghi *Importer) GetHomePageUrl() string {
	ghi.initClient()
	return ghi.repository.GetHomepage()
}
