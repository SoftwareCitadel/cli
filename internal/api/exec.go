package api

import (
	"bytes"
	"citadel/internal/util"
	"fmt"
	"io"
	"net/http"
)

func ExecuteCommand(
	organizationSlug string,
	projectSlug string,
	applicationSlug string,
	command string,
) error {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil
	}

	url := RetrieveApiBaseUrl() + "/organizations/" + organizationSlug + "/projects/" + projectSlug + "/applications/" + applicationSlug + "/exec"

	payload := []byte(`{"command": "` + command + `"}`)
	body := bytes.NewBuffer(payload)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/text")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Println("No running container found.")
		return nil
	}

	if resp.StatusCode != 200 {
		fmt.Println("An error occurred while executing the command.")
		return nil
	}

	// Print the output of the command to the console (response body)
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	fmt.Print(string(responseData))

	return nil
}
