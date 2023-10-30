package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/softwarecitadel/cli/internal/api"
)

func newChooseProjectPromptModel() SelectModel {
	projects, err := api.RetrieveProjects()
	if err != nil {
		fmt.Println("Failed to retrieve projects")
		os.Exit(1)
	}

	choices := []SelectChoice{}
	for _, project := range projects {
		choices = append(choices, SelectChoice{
			Name: project.Name,
			ID:   project.ID,
			Slug: project.Slug,
		})
	}

	return NewSelectModel("Which project would you like to deploy to?", choices)
}

func SelectProject() string {
	m := newChooseProjectPromptModel()
	res, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return res.(SelectModel).Choice.Slug
}
