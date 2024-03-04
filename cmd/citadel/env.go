package citadel

import (
	"fmt"
	"os"

	"citadel/internal/api"
	"citadel/internal/cli"
	"citadel/internal/tui"
	"citadel/internal/util"

	"github.com/charmbracelet/bubbles/table"
	"github.com/spf13/cobra"
)

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		applicationId, err := util.RetrieveApplicationIdFromProjectConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		envs, err := api.RetrieveEnvironmentVariables(applicationId)
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
		applicationId, err := util.RetrieveApplicationIdFromProjectConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(args) == 0 {
			fmt.Println("Please provide a least a pair of key and value")
			os.Exit(1)
		}

		showRedeployChoice, err := api.SetEnvironmentVariables(applicationId, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Environment variable set.")

		if !showRedeployChoice {
			return
		}

		shouldRedeploy := cli.AskYesOrNo("Do you want to (re)deploy your application?")
		if !shouldRedeploy {
			return
		}

		err = api.RedeployApplication(applicationId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Redeploying application...")
	},
}

var envLoadCmd = &cobra.Command{
	Use:     "load",
	Short:   "Load environment variables from file",
	Example: "citadel load .env",
	Run: func(cmd *cobra.Command, args []string) {
		applicationId, err := util.RetrieveApplicationIdFromProjectConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(args) != 1 {
			fmt.Println("Please provide the path to the file from which you want to load the environment variables")
			os.Exit(1)
		}

		vars, err := util.RetrieveEnvironmentVariablesFromFile(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		showRedeployChoice, err := api.SetEnvironmentVariables(applicationId, vars)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Environment variables loaded.")

		if !showRedeployChoice {
			return
		}

		shouldRedeploy := cli.AskYesOrNo("Do you want to (re)deploy your application?")
		if !shouldRedeploy {
			return
		}

		err = api.RedeployApplication(applicationId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Redeploying application...")
	},
}
