package server

import (
	"server/filefunc"
	"server/srv_conf"
	"server/srv_sec"

	"log"
)

func ServerInit(app_dir string) error {

	srv_yaml := app_dir + "/srv.yaml"
	if !filefunc.IsExists(srv_yaml) {
		// write the srv_conf file
		err := srv_conf.WriteConfigFile(app_dir)
		if err != nil {
			return err
		}
	}
	srv_conf.Srv_ReadConfig()
	srv_conf.SetPaths()

	// Check for .crt/.key files
	srv_sec.CheckTLS(app_dir, srv_conf.TLSKeySize())

	// Initialize JWT secret
	srv_sec.Env_SetSecret()

	// Check for static folder
	CheckFolder()

	appdb := srv_conf.GetString("app_db")
	if !filefunc.IsExists(appdb) {
		appdb = srv_conf.DataDir + "/app.db"
		log.Println("Creating", appdb)
		filefunc.CreateFile(appdb)
	}

	return nil

}

func CheckFolder() error {

	_stat_dir := srv_conf.StaticDir
	if filefunc.IsExists(_stat_dir) {
		log.Println("Deleting", _stat_dir)
		err := filefunc.DeleteFolder_FR(_stat_dir)
		if err != nil {
			return err
		}
	}

	err := filefunc.CreateFolder(_stat_dir)
	if err != nil {
		return err
	}

	assets_dir := srv_conf.AssetsDir
	if !filefunc.IsExists(assets_dir) {
		filefunc.CreateFolder(assets_dir)
		filefunc.CreateFolder(srv_conf.DataDir)
		filefunc.CreateFolder(srv_conf.ReportsDir)

		log.Println("Created assed dir:", assets_dir)
	}

	return nil
}
