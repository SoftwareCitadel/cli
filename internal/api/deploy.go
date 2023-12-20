package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"citadel/internal/util"

	"github.com/alevinval/sse/pkg/eventsource"
)

func DeployFromTarball(tarball io.ReadCloser, projectSlug string, applicationSlug string) (bool, error) {
	// Retrieve the token from the config file
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return false, err
	}

	// Create a new HTTP request
	url := ApiBaseUrl + "/projects/" + projectSlug + "/applications/" + applicationSlug + "/deploy"

	// Create the request, sending the gzipBuf as a form data field named "tarball"
	form := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(form)
	part, err := writer.CreateFormFile("tarball", "tarball")
	if err != nil {
		return false, err
	}
	_, err = io.Copy(part, tarball)
	if err != nil {
		return false, err
	}
	err = writer.Close()
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", url, form)
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the HTTP request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	var response struct {
		Healtcheck bool `json:"healtcheck"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, err
	}

	return response.Healtcheck, nil
}

func ShowBuildLogs(
	projectSlug string,
	applicationSlug string,
) {
	baseURL := ApiBaseUrl
	url := baseURL + "/projects/" + projectSlug + "/applications/" + applicationSlug + "/logs/stream?scope=builder"

	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	es, err := eventsource.New(url, eventsource.WithBearerTokenAuth(token))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		select {
		case event := <-es.MessageEvents():

			if strings.Contains(event.Data, "Main child exited") {
				if strings.Contains(event.Data, "Main child exited normally with code: 1") {
					fmt.Println("🔴 Build failed.")
					os.Exit(1)
				} else {
					fmt.Println("\n🚀 Build succeeded. Deploying application...\n")
					fmt.Println("\n🔗 Monitor the deployment at https://console.softwarecitadel.com/projects/" + projectSlug + "/applications/" + applicationSlug + "/logs\n")
					return
				}
			}

			fmt.Println(event.Data)
		}
	}
}

func RedeployApplication(
	projectSlug string,
	applicationSlug string,
) error {
	// Retrieve the token from the config file
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return err
	}

	// Create a new HTTP request
	url := ApiBaseUrl + "/api/projects/" + projectSlug + "/applications/" + applicationSlug + "/redeploy"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the HTTP request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	return nil
}
