package tui

import (
	"citadel/internal/api"
	"citadel/internal/util"
	"errors"
	"fmt"
	"os"

	"github.com/sveltinio/prompti/input"
)

func CreateProject() string {
	questionPrompt := &input.Config{
		Message:     "What's the name of your project?",
		Placeholder: "acme-inc",
		ErrorMsg:    "Please enter a valid project name",
		ValidateFunc: func(s string) error {
			if len(s) < 3 {
				return errors.New("Project name must be at least 3 characters")
			}
			if len(s) > 50 {
				return errors.New("Project name must be at most 50 characters")
			}
			if !util.IsSlug(s) {
				return errors.New("Project name must be a valid slug")
			}
			return nil
		},
	}

	projectName, _ := input.Run(questionPrompt)

	project, err := api.CreateProject(projectName)
	if err != nil {
		fmt.Println("An error occurred while creating the project.")
		os.Exit(1)
	}

	return project.Slug
}
