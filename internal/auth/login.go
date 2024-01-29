package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"citadel/internal/api"
)

type AuthenticationSessionResponse struct {
	SessionId string `json:"sessionId"`
}

type WaitForLoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

func GetAuthenticationSessionId() (string, error) {
	resp, err := http.Get(api.RetrieveApiBaseUrl() + "/auth/cli/session")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var auth AuthenticationSessionResponse

	err = json.Unmarshal(body, &auth)
	if err != nil {
		return "", err
	}

	return auth.SessionId, nil
}

func WaitForLogin(sessionId string) (string, error) {
	resp, err := http.Get(api.RetrieveApiBaseUrl() + "/auth/cli/" + sessionId + "/wait")
	if err != nil {
		fmt.Println("Whoops. There was an error while trying to get an authentication session.")
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var wait WaitForLoginResponse

	err = json.Unmarshal(body, &wait)
	if err != nil {
		return "", err
	}

	if wait.Status == "pending" {
		time.Sleep(1 * time.Second)
		return WaitForLogin(sessionId)
	}

	return wait.Token, nil
}
