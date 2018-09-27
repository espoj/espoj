package gitlab

import (
	"github.com/samitc/espoj/espoj/data"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type Importer struct {
	Url       string
	Token     string
	client    *gitlab.Client
	projectId int
	project   *gitlab.Project
	commits   []*gitlab.Commit
}

func getProjectNamespace(url string) string {
	lastSlash := strings.LastIndex(url, "/")
	lastSlash = strings.LastIndex(url[:lastSlash-1], "/")
	return url[lastSlash+1:]
}
func (gl *Importer) readCommits() {
	const perPage = 100
	trueVal := true
	commits, resp, _ := gl.client.Commits.ListCommits(gl.projectId, &gitlab.ListCommitsOptions{All: &trueVal, ListOptions: gitlab.ListOptions{PerPage: perPage}})
	var pCommits []*gitlab.Commit
	for resp.NextPage != 0 {
		pCommits, resp, _ = gl.client.Commits.ListCommits(gl.projectId, &gitlab.ListCommitsOptions{All: &trueVal, ListOptions: gitlab.ListOptions{PerPage: resp.ItemsPerPage, Page: resp.NextPage}})
		commits = append(commits, pCommits...)
	}
	gl.commits = commits
}
func (gl *Importer) initClient() {
	if gl.client == nil {
		gl.client = gitlab.NewClient(nil, gl.Token)
		gl.project, _, _ = gl.client.Projects.GetProject(getProjectNamespace(gl.Url))
		gl.projectId = gl.project.ID
		gl.readCommits()
	}
}
func (gl *Importer) GetUrl() string {
	gl.initClient()
	return gl.project.WebURL
}
func (gl *Importer) GetName() string {
	gl.initClient()
	return gl.project.Name
}
func (gl *Importer) GetDescription() string {
	gl.initClient()
	return gl.project.Description
}
func (gl *Importer) GetCommits() []data.Commit {
	gl.initClient()
	baseUrl := gl.GetUrl() + "/commit/"
	commits := make([]data.Commit, 0, len(gl.commits))
	for _, com := range gl.commits {
		commits = append(commits, data.Commit{Sha: com.ID, ShowUrl: baseUrl + com.ID, ParentSha: &com.ParentIDs, Time: com.AuthoredDate,
			Message: &com.Message})
	}
	return commits
}
func (gl *Importer) GetCloneUrl() string {
	gl.initClient()
	return gl.project.HTTPURLToRepo
}
func (gl *Importer) GetHomePageUrl() string {
	return ""
}
