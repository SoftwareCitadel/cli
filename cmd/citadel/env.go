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
	Short: "List environments",
	Run: func(cmd *cobra.Command, args []string) {
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

		envs, err := api.RetrieveEnvironmentVariables(projectSlug, applicationSlug)
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

		// Check if there is at least one argument
		if len(args) > 1 {
			fmt.Println("Please provide key and value")
			os.Exit(1)
		}

		showRedeployChoice, err := api.SetEnvironmentVariable(projectSlug, applicationSlug, args)
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

		err = api.RedeployApplication(projectSlug, applicationSlug)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Redeploying application...")
	},
}
