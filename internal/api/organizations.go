package api

import (
	"encoding/json"
	"net/http"
	"citadel/internal/util"
)

type Organization struct {
	ID   string	`json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func RetrieveOrganizations() ([]Organization, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil, err
	}

	url := RetrieveApiBaseUrl() + "/organizations"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	var organizations []Organization
	err = json.NewDecoder(resp.Body).Decode(&organizations)
	if err != nil {
		return nil, err
	}

	return organizations, nil
}
