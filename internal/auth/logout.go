package auth

import "github.com/softwarecitadel/cli/internal/util"

func Logout() error {
	if err := util.RemoveConfigFile(); err != nil {
		return err
	}

	return nil
}
