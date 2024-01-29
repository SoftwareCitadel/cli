package citadel

import (
	"fmt"
	"os"

	"citadel/internal/api"
	"citadel/internal/auth"
	"citadel/internal/tui"
	"citadel/internal/util"

	bspinner "github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/sveltinio/prompti/input"
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

		return
	}

	consoleUrl, _ := cmd.Flags().GetString("console-url")
	if consoleUrl != "" {
		err := util.StoreConsoleUrl(consoleUrl)
		if err != nil {
			fmt.Println("Whoops. There was an error while trying to store your console URL.")
			os.Exit(1)
		}
	} else {
		err := askAndStoreConsoleUrl()
		if err != nil {
			fmt.Println("Whoops. There was an error while trying to store your console URL.")
			os.Exit(1)
		}
	}

	s := bspinner.New()
	s.Spinner = bspinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#0099FF"))
	interrupt, err := tui.Run(s, runLoginTUI)
	if err != nil {
		fmt.Println("Whoops. There was an error while trying to log you in.")
	}

	if !interrupt {
		fmt.Println("\n🚀 You are now logged in to Software Citadel!")
	}
}

func runLoginTUI(msg tui.BasicSpinnerMessager) error {
	msg.SetStatus("Requesting authentication session")

	sessionId, err := auth.GetAuthenticationSessionId()
	if err != nil {
		return err
	}

	url := api.RetrieveApiBaseUrl() + "/auth/cli/" + sessionId
	fmt.Printf("\nOpening browser to %s\n", url)

	msg.SetStatus("Opening browser...")

	err = util.OpenInBrowser(api.RetrieveApiBaseUrl() + "/auth/cli/" + sessionId)
	if err != nil {
		return err
	}

	msg.SetStatus("Waiting for authentication...")

	token, err := auth.WaitForLogin(sessionId)
	if err != nil {
		return err
	}

	msg.SetStatus("Storing authentication token...")
	err = util.StoreJWTToken(token)
	if err != nil {
		return err
	}

	return nil
}

func askAndStoreConsoleUrl() error {
	consoleUrl, err := input.Run(&input.Config{
		Message:      "What's the console URL?",
		Placeholder:  "https://console.softwarecitadel.com",
		ErrorMsg:     "Please enter a valid console URL",
		ValidateFunc: util.UrlValidateFunc,
	})
	if err != nil {
		return err
	}

	err = util.StoreConsoleUrl(consoleUrl)
	if err != nil {
		return err
	}

	return nil
}
