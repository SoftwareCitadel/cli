package tui

import (
	"citadel/internal/api"
	"citadel/internal/util"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sveltinio/prompti/input"
)

func CreateApplication(
	orgSlug string,
	projectSlug string,
) string {
	questionPrompt := &input.Config{
		Message:      "What's the name of your application?",
		Placeholder:  "webapp",
		ErrorMsg:     "Please enter a valid application name",
		ValidateFunc: util.SlugValidateFunc,
	}

	applicationName, err := input.Run(questionPrompt)
	if err != nil {
		fmt.Println("An error occurred while creating the application.")
		os.Exit(1)
	}

	m := newChooseComputingSpecs()
	res, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	choice := res.(SelectModel).Choice
	splittedChoice := strings.Split(choice.ID, "x-")
	cpu := splittedChoice[0] + "x"
	memory := splittedChoice[1]

	application, err := api.CreateApplication(orgSlug, projectSlug, applicationName, cpu, memory)
	if err != nil {
		fmt.Println("\n🔴 " + err.Error())
		os.Exit(1)
	}

	return application.Slug
}

func newChooseComputingSpecs() SelectModel {
	choices := []SelectChoice{
		{
			Name: "shared-cpu-1x 256MB - $3.88/mo",
			ID:   "shared-cpu-1x-256MB",
		},
		{
			Name: "shared-cpu-1x 512MB - $6.38/mo",
			ID:   "shared-cpu-1x-512MB",
		},
		{
			Name: "shared-cpu-1x 1GB - $11.4/mo",
			ID:   "shared-cpu-1x-1GB",
		},
		{
			Name: "shared-cpu-1x 2GB - $21.4/mo",
			ID:   "shared-cpu-1x-2GB",
		},
		{
			Name: "shared-cpu-2x 512MB - $7.78/mo",
			ID:   "shared-cpu-2x-512MB",
		},
		{
			Name: "shared-cpu-2x 1GB - $15.28/mo",
			ID:   "shared-cpu-2x-1GB",
		},
		{
			Name: "shared-cpu-2x 2GB - $25.28/mo",
			ID:   "shared-cpu-2x-2GB",
		},
		{
			Name: "shared-cpu-2x 4GB - $45.3/mo",
			ID:   "shared-cpu-2x-4GB",
		},
		{
			Name: "shared-cpu-4x 1GB - $15.56/mo",
			ID:   "shared-cpu-4x-1GB",
		},
		{
			Name: "shared-cpu-4x 2GB - $33.06/mo",
			ID:   "shared-cpu-4x-2GB",
		},
		{
			Name: "shared-cpu-4x 4GB - $53.08/mo",
			ID:   "shared-cpu-4x-4GB",
		},
		{
			Name: "shared-cpu-4x 8GB - $93.1/mo",
			ID:   "shared-cpu-4x-8GB",
		},
		{
			Name: "shared-cpu-8x 2GB - $31.1/mo",
			ID:   "shared-cpu-8x-2GB",
		},
		{
			Name: "shared-cpu-8x 4GB - $68.62/mo",
			ID:   "shared-cpu-8x-4GB",
		},
		{
			Name: "performance-cpu-1x 2GB - $62.0/mo",
			ID:   "performance-cpu-1x-2GB",
		},
		{
			Name: "performance-cpu-1x 4GB - $82.02/mo",
			ID:   "performance-cpu-1x-4GB",
		},
		{
			Name: "performance-cpu-1x 8GB - $122.04/mo",
			ID:   "performance-cpu-1x-8GB",
		},
		{
			Name: "performance-cpu-2x 4GB - $124.0/mo",
			ID:   "performance-cpu-2x-4GB",
		},
		{
			Name: "performance-cpu-2x 8GB - $184.04/mo",
			ID:   "performance-cpu-2x-8GB",
		},
		{
			Name: "performance-cpu-2x 16GB - $264.08/mo",
			ID:   "performance-cpu-2x-16GB",
		},
		{
			Name: "performance-cpu-4x 8GB - $248.0/mo",
			ID:   "performance-cpu-4x-8GB",
		},
		{
			Name: "performance-cpu-4x 16GB - $388.08/mo",
			ID:   "performance-cpu-4x-16GB",
		},
		{
			Name: "performance-cpu-4x 32GB - $548.16/mo",
			ID:   "performance-cpu-4x-32GB",
		},
		{
			Name: "performance-cpu-8x 16GB - $496.0/mo",
			ID:   "performance-cpu-8x-16GB",
		},
		{
			Name: "performance-cpu-8x 32GB - $796.16/mo",
			ID:   "performance-cpu-8x-32GB",
		},
		{
			Name: "performance-cpu-8x 64GB - $1116.32/mo",
			ID:   "performance-cpu-8x-64GB",
		},
		{
			Name: "performance-cpu-16x 32GB - $992.02/mo",
			ID:   "performance-cpu-16x-32GB",
		},
		{
			Name: "performance-cpu-16x 64GB - $1612.32/mo",
			ID:   "performance-cpu-16x-64GB",
		},
		{
			Name: "performance-cpu-16x 128GB - $2248.64/mo",
			ID:   "performance-cpu-16x-128GB",
		},
	}

	return NewSelectModel("Which computing specs would you like to use?", choices)
}
