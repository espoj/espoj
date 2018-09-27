package repository

import "github.com/samitc/espoj/espoj/data"

type Importer interface {
	GetUrl() string
	GetName() string
	GetDescription() string
	GetCommits() []data.Commit
	GetCloneUrl() string
	GetHomePageUrl() string
}

func CreateRepository(repositoryImport Importer) data.Repository {
	return data.Repository{repositoryImport.GetUrl(), repositoryImport.GetName(), repositoryImport.GetDescription(),repositoryImport.GetCloneUrl() ,repositoryImport.GetCommits()}
}
