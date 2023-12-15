package tui

import (
	"fmt"
	"os"

	"citadel/internal/api"
)

func newChooseProjectPromptModel(title string) SelectModel {
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
	choices = append(choices, SelectChoice{
		Name: "Create a new project",
		ID:   "",
		Slug: "",
	})

	return NewSelectModel(title, choices)
}

func SelectProject(
	title string,
) string {
	m := newChooseProjectPromptModel(title)

	choice, err := m.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return choice.Slug
}
