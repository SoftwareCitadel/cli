package tui

import (
	"citadel/internal/api"
	"citadel/internal/util"
	"fmt"
	"os"
	"strings"

	"github.com/alevinval/sse/pkg/eventsource"
	tea "github.com/charmbracelet/bubbletea"
)

func StreamBuildLogs(
	projectSlug string,
	applicationSlug string,
) {
	streamModel := NewStreamModel("Building...")

	buildFailed := false

	streamModel.Run(func(p *tea.Program) {
		baseURL := api.RetrieveApiBaseUrl()
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
				if strings.Contains(event.Data, "Main child exited") || strings.Contains(event.Data, "Pushed Docker image, built with Nixpacks.") || strings.Contains(event.Data, "Pushed Docker image.") {
					if strings.Contains(event.Data, "Main child exited normally with code: 1") {
						buildFailed = true
					}
					p.Quit()
				}
				splitted := strings.Split(event.Data, " | ")
				subtitle := splitted[0]
				title := splitted[1]
				p.Send(StreamModelResultMsg{
					Title:    title,
					Subtitle: subtitle,
				})
			}
		}
	})

	if buildFailed {
		fmt.Println("🔴 Build failed.")
		os.Exit(1)
	} else {
		fmt.Println("🚀 Build succeeded. Deploying application...")
		fmt.Println("\n🔗 Monitor the deployment at " + api.RetrieveApiBaseUrl() + "/projects/" + projectSlug + "/applications/" + applicationSlug + "/logs\n")
	}
}
