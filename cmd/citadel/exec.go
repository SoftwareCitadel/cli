package citadel

import (
	"citadel/internal/api"
	"citadel/internal/auth"
	"citadel/internal/util"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [command]",
	Run:   runExec,
	Short: "Execute a command on the running container",
}

func runExec(cmd *cobra.Command, args []string) {
	if !auth.IsLoggedIn() {
		fmt.Println("You are not logged in to Software Citadel.")
		os.Exit(1)
	}

	projectSlug, err := util.RetrieveProjectSlugFromProjectConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	applicationSlug, err := util.RetrieveApplicationSlugFromProjectConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) != 1 {
		fmt.Println("Please provide a command to execute.")
		os.Exit(1)
	}

	command := args[0]

	err = api.ExecuteCommand(projectSlug, applicationSlug, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
