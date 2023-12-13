package citadel

import (
	"fmt"
	"os"

	"citadel/internal/auth"
	"citadel/internal/cli"
	"citadel/internal/tui"
	"citadel/internal/util"

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

	if projectSlug == "" {
		projectSlug = tui.SelectProject()
		if projectSlug == "" {
			projectSlug = tui.CreateProject()
		}
	}

	if applicationSlug == "" {
		applicationSlug = tui.SelectApplication(projectSlug)
		if applicationSlug == "" {
			applicationSlug = tui.CreateApplication(projectSlug)
		}
	}

	err := util.InitializeConfigFile(projectSlug, applicationSlug)
	if err != nil {
		fmt.Println("Failed to initialize Software Citadel project.")
		return
	}

	fmt.Println("Project initialized.")
}
