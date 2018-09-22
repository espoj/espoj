package appveyor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

var httpClient = http.DefaultClient

const baseUrl = "https://ci.appveyor.com/api/"

type Importer struct {
	AccountName string
	Token       string
}

func (appveyor *Importer) getRes(restReq string) string {
	reqS := baseUrl + restReq
	req, _ := http.NewRequest("GET", reqS, nil)
	req.Header.Set("Authorization", "Bearer "+appveyor.Token)
	req.Header.Set("Content-type", "application/json")
	res, _ := httpClient.Do(req)
	defer func() { res.Body.Close() }()
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	return string(bodyBytes)
}
func (appveyor *Importer) getUsers() string {
	restReq := "users"
	return appveyor.getRes(restReq)
}
func (appveyor *Importer) getProjects() string {
	restReq := "projects"
	return appveyor.getRes(restReq)
}
func (appveyor *Importer) getHistory(projectSlug string) string {
	//GET /api/projects/{accountName}/{projectSlug}/history?recordsNumber={records-per-page}[&startBuildId={buildId}&branch={branch}]
	restReq := "projects/" + appveyor.AccountName + "/" + projectSlug + "/history?recordsNumber=100000"
	return appveyor.getRes(restReq)
}
func (appveyor *Importer) getBuild(projectSlug, buildNumber string) string {
	//GET /api/projects/{accountName}/{projectSlug}/build/{buildVersion}
	restReq := "projects/" + appveyor.AccountName + "/" + projectSlug + "/build/" + buildNumber
	return appveyor.getRes(restReq)
}
func (appveyor *Importer) getArtifacts(jobId string) string {
	restReq := "buildjobs/" + jobId + "/artifacts"
	return appveyor.getRes(restReq)
}
func (appveyor *Importer) downloadArtifact(jobId, artifact, prefixDirectory string) {
	//$apiUrl/buildjobs/$jobId/artifacts/$artifactFileName
	restReq := "buildjobs/" + jobId + "/artifacts/" + artifact
	reqS := baseUrl + restReq
	req, _ := http.NewRequest("GET", reqS, nil)
	req.Header.Set("Authorization", "Bearer "+appveyor.Token)
	res, _ := httpClient.Do(req)
	defer func() {
		if res != nil {
			res.Body.Close()
		} else {
			fmt.Println("res null" + jobId + "," + artifact)
		}
	}()
	if res != nil {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		ioutil.WriteFile(prefixDirectory+artifact, bodyBytes, 0644)
	} else {
		appveyor.downloadArtifact(jobId, artifact, prefixDirectory)
	}
}
func (appveyor *Importer) DownloadAllArtifacts(outputFolder string) error {
	if outputFolder[len(outputFolder)-1:] != string(os.PathSeparator) {
		outputFolder += string(os.PathSeparator)
	}
	var wg sync.WaitGroup
	slugs := appveyor.getAllSlugs()
	for _, slug := range slugs {
		hist := appveyor.getAllHistory(slug)
		for _, his := range hist {
			jobs := appveyor.getAllJobsIds(slug, his)
			for _, job := range jobs {
				folder := outputFolder + job + `\`
				err:=os.Mkdir(folder, 0777)
				if err!=nil{
					return err
				}
				artifacts := appveyor.getAllArtifacts(job)
				wg.Add(len(artifacts))
				for _, artifact := range artifacts {
					go func(jobId, artifact, folderPath string) {
						defer wg.Done()
						appveyor.downloadArtifact(jobId, artifact, folderPath)
					}(job, artifact, folder)
				}
			}
		}
	}
	wg.Wait()
	return nil
}
func (appveyor *Importer) getAllHistory(slug string) []string {
	return getAllFromJson(appveyor.getHistory(slug), "version")
}
func (appveyor *Importer) getAllSlugs() []string {
	return getAllFromJson(appveyor.getProjects(), "slug")
}
func (appveyor *Importer) getAllJobsIds(slug, buildId string) []string {
	return getAllFromJson(appveyor.getBuild(slug, buildId), "jobId")
}
func (appveyor *Importer) getAllArtifacts(jobId string) []string {
	return getAllFromJson(appveyor.getArtifacts(jobId), "fileName")
}
