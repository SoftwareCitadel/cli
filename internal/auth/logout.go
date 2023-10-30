package auth

import "citadel/internal/util"

func Logout() error {
	if err := util.RemoveConfigFile(); err != nil {
		return err
	}

	return nil
}
