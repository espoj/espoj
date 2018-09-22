package data

type Commit struct {
	Sha     string
	parentSha string
	ShowUrl string
}
type Issue struct {
}
type PullRequest struct {
}
type Repository struct {
	Url, Name, Description string
	Commits                []Commit
}