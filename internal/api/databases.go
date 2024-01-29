package api

import (
	"bytes"
	"citadel/internal/util"
	"encoding/json"
	"errors"
	"net/http"
)

func CreateDatabase(
	projectSlug string,
	dbms string,
	databaseName string,
	databaseUsername string,
	databasePassword string,
	diskSize string,
) (string, string, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return "", "", nil
	}

	url := RetrieveApiBaseUrl() + "/projects/" + projectSlug + "/databases"
	payload := `{"name": "` + databaseName + `",`
	payload += `"dbms": "` + dbms + `",`
	payload += `"username": "` + databaseUsername + `",`
	if databasePassword != "" {
		payload += `"password": "` + databasePassword + `",`
	}
	payload += `"diskSize": ` + diskSize + `}`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var response struct {
		ConnectionString string `json:"connectionString"`
		Slug             string `json:"slug"`
	}

	if resp.StatusCode == 400 {
		var output map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&output)
		if err != nil {
			return "", "", err
		}
		if output["message"] != nil {
			return "", "", errors.New(output["message"].(string))
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", "", err
	}
	return response.ConnectionString, response.Slug, nil
}
