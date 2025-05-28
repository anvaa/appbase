package srv_int

import (
	"srv/filefunc"
	"srv/server"
	"srv/srv_conf"

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
	server.CheckTLS(app_dir, srv_conf.GetInt("tls_keysize"))

	// Check for static folder
	CheckFolder()

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
		filefunc.CreateFolder(srv_conf.QRImgDir)
		filefunc.CreateFolder(srv_conf.BarImgDir)
		log.Println("Created assed dir:", assets_dir)
	}

	return nil
}
