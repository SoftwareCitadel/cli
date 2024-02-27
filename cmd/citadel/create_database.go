package citadel

import (
	"citadel/internal/auth"
	"citadel/internal/tui"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var createDatabaseCmd = &cobra.Command{
	Use:   "databases create",
	Run:   runCreateDatabase,
	Short: "Create a database",
}

func init() {
	rootCmd.AddCommand(createDatabaseCmd)
}

func runCreateDatabase(cmd *cobra.Command, args []string) {
	if !auth.IsLoggedIn() {
		fmt.Println("You must be logged in to initialize a project.\nPlease run `citadel auth login` to log in.")
		os.Exit(1)
	}

	orgSlug := tui.SelectOrganization()

	projectSlug := tui.SelectProject(orgSlug, "Select a project to create a database for")
	if projectSlug == "" {
		projectSlug = tui.CreateProject(orgSlug)
	}

	tui.CreateDatabase(orgSlug, projectSlug)
}
