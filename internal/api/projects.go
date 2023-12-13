package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"citadel/internal/util"
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

func CreateProject(projectName string) (Project, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return Project{}, err
	}

	url := ApiBaseUrl + "/projects"

	payload := []byte(`{"name": "` + projectName + `"}`)
	body := bytes.NewBuffer(payload)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return Project{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Project{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Project{}, errors.New("An error occurred while creating the project.")
	}

	var project Project
	err = json.NewDecoder(resp.Body).Decode(&project)
	if err != nil {
		return Project{}, err
	}

	return project, nil
}
