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
	initCmd.Flags().StringP("application-id", "a", "", "Application id to use for initialization (optional)")

	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	applicationId, _ := cmd.Flags().GetString("application-id")

	if !auth.IsLoggedIn() {
		fmt.Println("You must be logged in to initialize a project.\nPlease run `citadel auth login` to log in.")
		os.Exit(1)
	}

	if util.IsAlreadyInitialized() {
		if !cli.AskYesOrNo("Software Citadel is already initialized. Do you want to reinitialize it?") {
			return
		}
	}

	// Retrieve or create application
	if applicationId == "" {
		applicationId = tui.SelectApplication()
		if applicationId == "" {
			applicationId = tui.CreateApplication()
		}
	}

	err := util.InitializeConfigFile(applicationId)
	if err != nil {
		fmt.Println("Failed to initialize Software Citadel project.")
		return
	}

	fmt.Println("\nCongratulations! Your project is now initialized.")
}
