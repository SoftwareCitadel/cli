package tui

import (
	"fmt"
	"os"

	"citadel/internal/api"
)

func newChooseOrganizationPromptModel() SelectModel {
	organizations, err := api.RetrieveOrganizations()
	if err != nil {
		fmt.Println("Failed to retrieve organizations", err)
		os.Exit(1)
	}
	fmt.Println(organizations)

	choices := []SelectChoice{}
	for _, organization := range organizations {
		choices = append(choices, SelectChoice{
			Name: organization.Name,
			ID:   organization.ID,
			Slug: organization.Slug,
		})
	}

	return NewSelectModel("Select an organization", choices)
}

func SelectOrganization() string {
	m := newChooseOrganizationPromptModel()

	choice, err := m.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return choice.Slug
}
