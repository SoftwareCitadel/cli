package citadel

import (
	"fmt"
	"os"

	"citadel/internal/api"
	"citadel/internal/cli"
	"citadel/internal/util"
	"github.com/spf13/cobra"
)

var postgresCmd = &cobra.Command{
	Use: "postgres",
}

var createPostgresCmd = &cobra.Command{
	Use: "create",
	Run: runCreatePostgres,
}

func init() {
	postgresCmd.AddCommand(createPostgresCmd)
	rootCmd.AddCommand(postgresCmd)
}

func runCreatePostgres(cmd *cobra.Command, args []string) {
	projectId, err := util.RetrieveProjectIdFromProjectConfig()
	if err != nil {
		fmt.Println("Failed to retrieve project id")
		os.Exit(1)
	}

	applicationName := cli.Ask("What is the name of the application you want to create to this postgres instance?")

	databaseUri, err := api.CreatePostgresToProject(projectId, applicationName)
	if err != nil {
		fmt.Println("Failed to create postgres instance to project " + projectId)
		os.Exit(1)
	}

	fmt.Println("Successfully Createed postgres instance to project " + projectId)
	fmt.Println("Your database URI is: " + databaseUri)
}
