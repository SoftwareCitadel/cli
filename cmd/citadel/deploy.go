package citadel

import (
	"fmt"
	"os"

	"citadel/internal/api"
	"citadel/internal/auth"
	"citadel/internal/util"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Run:   runDeploy,
	Short: "Deploy your project",
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	if !auth.IsLoggedIn() {
		fmt.Println("You must be logged in to deploy a project.")
		fmt.Println("Please run `citadel auth login` to log in.")
		return
	}

	if !util.IsAlreadyInitialized() {
		fmt.Println("Software Citadel is not initialized. Please run `citadel init` to initialize it.")
		return
	}

	if !checkDockerfileExists() {
		fmt.Println("Dockerfile not found. Please create a Dockerfile in the root of your project.")
		return
	}

	fmt.Println("Deploying...")

	tarball, err := util.MakeTarball()
	if err != nil {
		fmt.Println(err)
		return
	}

	projectId, err := util.RetrieveProjectIdFromProjectConfig()
	if err != nil {
		fmt.Println("Failed to retrieve project id")
		os.Exit(1)
	}

	applicationId, err := util.RetrieveApplicationIdFromProjectConfig()
	if err != nil {
		fmt.Println("Failed to retrieve application id")
		os.Exit(1)
	}

	err = api.DeployFromTarball(tarball, projectId, applicationId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api.ShowBuildLogs(projectId, applicationId)

	fmt.Println("Deployed!")
}

func checkDockerfileExists() bool {
	if _, err := os.Stat("./Dockerfile"); os.IsNotExist(err) {
		return false
	}

	return true
}
