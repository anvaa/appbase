package users

import (
	"srv/filefunc"
	"usr/user_embed"
)

// User_Init writes static html to disk
// and initializes the user configuration file usr.yaml
func UserInit(app_path string) error {

	// user embed files
	err := user_embed.User_EmbedFiles()
	if err != nil {
		return err
	}

	// Check/make usr.yaml
	if !filefunc.IsExists(app_path + "/usr.yaml") {
		WriteConfigFile(app_path)
	}
	ReadConfig()

	return nil
}
