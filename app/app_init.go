package app

import (
	"fmt"
	"log"

	"app/app_conf"
	"app/app_db"
	"app/app_embed"

	"server/filefunc"
	"server/srv_conf"

)

func AppInit(app_folder string) error {

	// write app static files to disk
	err := app_embed.App_EmbedFiles()
	if err != nil {
		return err
	}

	// check for app config file
	configFile := app_folder + "/app.yaml"
	if !filefunc.IsExists(configFile) {
		log.Println("No app.yaml file found. Creating", configFile)
		app_conf.WriteDefaultConfig(app_folder)
	}
	app_conf.ReadConfig() // read the config file

	

	// connect/sync to the app database
	app_db.CnnAppDB(srv_conf.GetDbConfig())

	// get default values fro DB
	getMenuDefaultId()

	return nil
}

func getMenuDefaultId() {
	// Get the next status ID
	sta_def := app_db.Sta_GetStaSubIDByType("[D]")
	for _, v := range sta_def {
		newid := v.ID
		_txt := fmt.Sprintf("sta_def%d", v.StatusID)
		app_conf.SetVal(_txt, newid)
	}

	typ_def := app_db.Typ_GetTypSubIDByType("[D]")
	for _, v := range typ_def {
		newid := v.ID
		_txt := fmt.Sprintf("typ_def%d", v.TypeID)
		app_conf.SetVal(_txt, newid)
	}

	_mnu := app_db.Mnu_GetMnuSubIDByType("[D]")
	for _, v := range _mnu {
		newid := v.ID
		_txt := fmt.Sprintf("mnu_def%d", v.MenuID)
		app_conf.SetVal(_txt, newid)
	}

}
