package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/softwarecitadel/cli/internal/util"
)

func RetrieveEnvironmentVariables(
	projectId string,
	applicationId string,
) (map[string]string, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil, err
	}

	url := ApiBaseUrl + "/api/projects/" + projectId + "/applications/" + applicationId + "/env"
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

	var response struct {
		Env map[string]string `json:"env"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Env, nil
}

func SetEnvironmentVariable(
	projectId string,
	applicationId string,
	args []string,
) (bool, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return false, err
	}

	data := "{"
	for i, arg := range args {
		// Check arg follows the pattern <key>=<value>
		if !strings.Contains(arg, "=") {
			return false, errors.New("Invalid argument: " + arg)
		}

		splittedArg := strings.Split(arg, "=")
		key := splittedArg[0]
		value := splittedArg[1]

		data += "\"" + key + "\":\"" + value + "\""
		if i < len(args)-1 {
			data += ","
		}
	}
	data += "}"

	body := bytes.NewBufferString(data)

	url := ApiBaseUrl + "/projects/" + projectId + "/applications/" + applicationId + "/env"
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return false, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, errors.New("Failed to set environment variable")
	}

	// Parse the response body { showRedeployChoice: true }
	var response struct {
		ShowRedeployChoice bool `json:"showRedeployChoice"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, err
	}

	return response.ShowRedeployChoice, nil
}
