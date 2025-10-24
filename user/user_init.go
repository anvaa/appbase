package user

import (
	"server/filefunc"
	"user/user_conf"
	"user/user_embed"
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
		user_conf.WriteConfigFile(app_path)
	}
	user_conf.ReadConfig()

	return nil
}
