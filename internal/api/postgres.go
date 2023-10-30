package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"citadel/internal/util"
)

// CreatePostgresToProject attaches a postgres instance to a project
// and returns the connection string
func CreatePostgresToProject(
	projectId string,
	applicationName string,
) (string, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(map[string]string{
		"name": applicationName,
	})

	url := ApiBaseUrl + "/api/projects/" + projectId + "/postgres"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Failed to attach postgres instance to project")
	}

	var response struct {
		DatabaseUri string `json:"databaseUri"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.DatabaseUri, nil
}
