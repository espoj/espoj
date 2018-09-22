package main

import (
	"github.com/samitc/espoj/espoj/data"
	"github.com/samitc/espoj/espoj/repository"
	"github.com/samitc/espoj/espoj/repository/github"
	"github.com/samitc/espoj/espoj/repository/gitlab"
	"net/http"
	"strings"
)

func createRepo(url string) (*data.Repository) {
	var repo data.Repository
	var success = false
	if strings.Contains(url, "github") {
		repo = repository.CreateRepository(&github.Importer{Url: url})
		success = true
	} else if strings.Contains(url, "gitlab") {
		repo = repository.CreateRepository(&gitlab.Importer{Url: url, Token: ""})
		success = true
	}
	if !success {
		return nil
	}
	return &repo
}
func main() {
	data1 := data.NewData()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		url := strings.TrimPrefix(request.URL.Path, "/")
		repo := data1.GetRepository(url[strings.LastIndex(url, "/")+1:])
		if repo == nil {
			repo = createRepo(url)
			data1.AddRepository(repo)
		}
		writer.Write([]byte("<html><head></head><body>"))
		if repo != nil {
			writer.Write([]byte("<h1><a href=" + repo.Url + ">" + repo.Name + "</a></h1><h2>commits</h2>"))
			for _, commit := range repo.Commits {
				writer.Write([]byte("<h5><a href=" + commit.ShowUrl + ">" + commit.Sha + "</a></h5>"))
			}
		}
		writer.Write([]byte("</body></html>"))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
