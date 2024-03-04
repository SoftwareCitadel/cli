package citadel

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"citadel/internal/util"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use: "open",
	Run: runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) {
	applicationId, err := util.RetrieveApplicationIdFromProjectConfig()
	if err != nil {
		fmt.Println("Failed to retrieve application id")
		os.Exit(1)
	}

	openBrowser("https://" + applicationId + ".softwarecitadel.app")
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
