package main

import (
	"citadel/cmd/citadel"
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/creativeprojects/go-selfupdate"
)

const (
	version = "dev"
)

func main() {
	err := update(version)
	if err != nil {
		fmt.Printf("error occurred while updating binary: %v\n", err)
		os.Exit(1)
	}

	citadel.Execute()
}

func update(version string) error {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("softwarecitadel/cli"))
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}
	if !found {
		return fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
	}

	if latest.Version() == version {
		fmt.Printf("Current version (%s) is the latest\n", version)
		return nil
	}

	fmt.Printf("Updating to version %s...\n", latest.Version())

	exe, err := os.Executable()
	if err != nil {
		return errors.New("could not locate executable path")
	}
	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	fmt.Printf("Successfully updated to version %s.\n", latest.Version())

	return nil
}
