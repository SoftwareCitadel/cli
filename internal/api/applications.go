package api

import (
	"bytes"
	"citadel/internal/util"
	"encoding/json"
	"errors"
	"net/http"
)

type Application struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func RetrieveApplications(organizationSlug string, projectSlug string) ([]Application, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil, err
	}

	url := RetrieveApiBaseUrl() + "/organizations/" + organizationSlug + "/projects/" + projectSlug + "/applications"
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

func CreateApplication(
	orgSlug string,
	projectSlug string,
	applicationName string,
	cpu string,
	memory string,
) (Application, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return Application{}, err
	}

	url := RetrieveApiBaseUrl() + "/organizations/" + orgSlug + "/projects/" + projectSlug + "/applications"
	payload := `{"name": "` + applicationName + `",`
	payload += `"cpu": "` + cpu + `",`
	payload += `"ram": "` + memory + `"}`

	body := bytes.NewBuffer([]byte(payload))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return Application{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Application{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		var output map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&output)
		if err != nil {
			return Application{}, err
		}
		if output["message"] != nil {
			return Application{}, errors.New(output["message"].(string))
		}
	}

	if resp.StatusCode != 200 {
		return Application{}, errors.New("An error occurred while creating the application.")
	}

	var application Application
	err = json.NewDecoder(resp.Body).Decode(&application)
	if err != nil {
		return Application{}, err
	}

	return application, nil
}
