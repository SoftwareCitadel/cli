package citadel

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/softwarecitadel/cli/internal/api"
	"github.com/softwarecitadel/cli/internal/cli"
	"github.com/softwarecitadel/cli/internal/tui"
	"github.com/softwarecitadel/cli/internal/util"
	"github.com/spf13/cobra"
)

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments",
	Run: func(cmd *cobra.Command, args []string) {
		projectId, err := util.RetrieveProjectIdFromProjectConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		applicationId, err := util.RetrieveApplicationIdFromProjectConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		envs, err := api.RetrieveEnvironmentVariables(projectId, applicationId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tui.ShowTable(
			[]table.Column{
				{Title: "Key", Width: 20},
				{Title: "Value", Width: 50},
			},
			envs,
			5,
		)
	},
}

var envSetCmd = &cobra.Command{
	Use:     "set <KEY>=<VALUE>",
	Short:   "Set environment",
	Example: "citadel set DATABASE_URL=postgresql://username:password@host:5432/mydb",
	Run: func(cmd *cobra.Command, args []string) {
		projectId, err := util.RetrieveProjectIdFromProjectConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		applicationId, err := util.RetrieveApplicationIdFromProjectConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Check if there is at least one argument
		if len(args) > 1 {
			fmt.Println("Please provide key and value")
			os.Exit(1)
		}

		showRedeployChoice, err := api.SetEnvironmentVariable(projectId, applicationId, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Environment variable set")

		if !showRedeployChoice {
			return
		}

		shouldRedeploy := cli.AskYesOrNo("Do you want to (re)deploy your application?")
		if !shouldRedeploy {
			return
		}

		err = api.RedeployApplication(projectId, applicationId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Redeploying application...")
	},
}
