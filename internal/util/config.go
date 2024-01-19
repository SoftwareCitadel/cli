package util

import (
	"os"

	"github.com/spf13/viper"
)

func RetrieveTokenFromConfig() (string, error) {
	configDir, err := InitConfigDir()
	if err != nil {
		return "", err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	token := viper.GetString("jwt")

	return token, nil
}

func RetrieveProjectSlugFromProjectConfig() (string, error) {
	viper.SetConfigName("citadel")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	projectSlug := viper.GetString("project_id")

	return projectSlug, nil
}

func RetrieveApplicationSlugFromProjectConfig() (string, error) {
	viper.SetConfigName("citadel")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	applicationSlug := viper.GetString("application_id")

	return applicationSlug, nil
}

func IsAlreadyInitialized() bool {
	vi := viper.New()
	vi.SetConfigName("citadel")
	vi.AddConfigPath(".")
	vi.SetConfigType("toml")

	err := vi.ReadInConfig()

	return err == nil
}

func RetrieveReleaseCommandFromProjectConfig() (string, error) {
	viper.SetConfigName("citadel")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	releaseCommand := viper.GetString("release_command")

	return releaseCommand, nil
}

func InitializeConfigFile(
	projectName string,
	applicationName string,
) error {
	vi := viper.New()
	vi.SetConfigName("citadel")
	vi.AddConfigPath(".")
	vi.SetConfigType("toml")

	vi.Set("project_id", projectName)
	vi.Set("application_id", applicationName)

	var fileExists bool

	_, err := os.Stat("citadel.toml")
	if err == nil {
		fileExists = true
	} else if os.IsNotExist(err) {
		fileExists = false
	} else {
		return err
	}

	if fileExists {
		err = vi.MergeInConfig()
		if err != nil {
			return err
		}
	} else {
		_, err = os.Create("citadel.toml")
		if err != nil {
			return err
		}
	}

	err = vi.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
