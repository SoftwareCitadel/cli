package citadel

import (
	"fmt"
	"os"

	"citadel/internal/auth"
	"citadel/internal/cli"
	"citadel/internal/tui"
	"citadel/internal/util"
	"citadel/internal/api"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Run:   runInit,
	Short: "Initialize a Software Citadel project",
}

func init() {
	initCmd.Flags().StringP("project", "p", "", "Project slug")
	initCmd.Flags().StringP("application", "a", "", "Application slug")

	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	organizationSlug, _ := cmd.Flags().GetString("organization")
	projectSlug, _ := cmd.Flags().GetString("project")
	applicationSlug, _ := cmd.Flags().GetString("application")

	if !auth.IsLoggedIn() {
		fmt.Println("You must be logged in to initialize a project.\nPlease run `citadel auth login` to log in.")
		os.Exit(1)
	}

	if util.IsAlreadyInitialized() {
		if !cli.AskYesOrNo("Software Citadel is already initialized. Do you want to reinitialize it?") {
			return
		}
	}

	if organizationSlug == "" {
		organizationSlug = tui.SelectOrganization()
	}

	if projectSlug == "" {
		projectSlug = tui.SelectProject(organizationSlug, "Which project would you like to deploy to?")
		if projectSlug == "" {
			projects, err := api.RetrieveProjects(organizationSlug)
			if err != nil {
				fmt.Println("Failed to retrieve projects")
				os.Exit(1)
			}
			fmt.Println(projects)
			// projectSlug = tui.CreateProject(organizationSlug)
		}
	}

	if applicationSlug == "" {
		applicationSlug = tui.SelectApplication(organizationSlug, projectSlug)
		if applicationSlug == "" {
			applicationSlug = tui.CreateApplication(organizationSlug, projectSlug)
		}
	}

	err := util.InitializeConfigFile(organizationSlug, projectSlug, applicationSlug)
	if err != nil {
		fmt.Println("Failed to initialize Software Citadel project.")
		return
	}

	fmt.Println("\nCongratulations! Your project is now initialized.")
}
