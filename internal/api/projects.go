package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/softwarecitadel/cli/internal/util"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ProjectWithApplications struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Applications []Application `json:"applications"`
}

type Application struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func RetrieveProjects() ([]Project, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil, err
	}

	url := ApiBaseUrl + "/projects"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return nil, err
	}

	var projects []Project
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func RetrieveApplications(projectId string) ([]Application, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil, err
	}

	url := ApiBaseUrl + "/projects/" + projectId + "/applications"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	var applications []Application
	err = json.NewDecoder(resp.Body).Decode(&applications)
	if err != nil {
		return nil, err
	}

	return applications, nil
}
