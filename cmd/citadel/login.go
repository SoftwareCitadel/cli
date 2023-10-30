package citadel

import (
	"fmt"
	"os"

	bspinner "github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/softwarecitadel/cli/internal/api"
	"github.com/softwarecitadel/cli/internal/auth"
	"github.com/softwarecitadel/cli/internal/tui"
	"github.com/softwarecitadel/cli/internal/util"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Run:   runLogin,
	Short: "Login to Software Citadel",
}

func runLogin(cmd *cobra.Command, args []string) {
	token, _ := cmd.Flags().GetString("token")
	if token != "" {
		err := util.StoreJWTToken(token)
		if err != nil {
			fmt.Println("Whoops. There was an error while trying to store your authentication token.")
			os.Exit(1)
		}

		fmt.Println("Authentication token stored")

		return
	}

	s := bspinner.New()
	s.Spinner = bspinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#0099FF"))
	interrupt, err := tui.Run(s, runLoginTUI)
	if err != nil {
		fmt.Println("Whoops. There was an error while trying to log you in.")
	}

	if !interrupt {
		fmt.Println("🚀 You are now logged in to Software Citadel!")
	}
}

func runLoginTUI(msg tui.BasicSpinnerMessager) error {
	msg.SetStatus("Requesting authentication session")

	sessionId, err := auth.GetAuthenticationSessionId()
	if err != nil {
		fmt.Println("Whoops. There was an error while trying to get an authentication session.")
		os.Exit(1)
	}

	url := api.ApiBaseUrl + "/auth/cli/" + sessionId
	fmt.Printf("Opening browser to %s\n", url)

	msg.SetStatus("Opening browser...")

	err = util.OpenInBrowser(api.ApiBaseUrl + "/auth/cli/" + sessionId)
	if err != nil {
		fmt.Println("Whoops. There was an error while trying to open your browser.")
		os.Exit(1)
	}

	msg.SetStatus("Waiting for authentication...")

	token, err := auth.WaitForLogin(sessionId)
	if err != nil {
		fmt.Println("Whoops. There was an error while trying to get your authentication token.")
		os.Exit(1)
	}

	msg.SetStatus("Storing authentication token...")
	err = util.StoreJWTToken(token)
	if err != nil {
		fmt.Println("Whoops. There was an error while trying to store your authentication token.")
		os.Exit(1)
	}

	return nil
}
