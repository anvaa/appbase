package app

import (
	"fmt"
	"log"

	"app/app_conf"
	"app/app_db"
	"app/app_embed"

	"srv/filefunc"
	"srv/srv_conf"
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

	appdb := app_conf.GetString("app_db")
	if !filefunc.IsExists(appdb) {
		appdb = srv_conf.DataDir + "/app.db"
		log.Println("Creating", appdb)
		filefunc.CreateFile(appdb)
	}

	// connect/sync to the app database
	app_db.CnnAppDB(appdb)

	// get default values fro DB
	setStatusNewId()
	setTypesNewId()

	return nil
}

func setStatusNewId() {
	// Get the next status ID
	sta_def := app_db.Sta_GetStaSubIDByType("[D]")
	for _, v := range sta_def {
		newid := v.ID
		_txt := fmt.Sprintf("stadef%d", v.StatusID)
		app_conf.SetVal(_txt, newid)
	}

	_mnu := app_db.Mnu_GetMnuSubIDByType("[D]")
	for _, v := range _mnu {
		newid := v.ID
		_txt := fmt.Sprintf("mnudef%d", v.MenuID)
		app_conf.SetVal(_txt, newid)
	}

}

func setTypesNewId() {
	// Get the next type ID
	typ_def := app_db.Typ_GetTypSubIDByType("[D]")
	for _, v := range typ_def {
		newid := v.ID
		_txt := fmt.Sprintf("typdef%d", v.TypeID)
		app_conf.SetVal(_txt, newid)
	}
}
