package app

import (

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

	return nil
}

func setStatusNewId() {
	// Get the next status ID
	newid := app_db.Sta_GetStatusIdByType("[D]")
	app_conf.SetVal("status_default", newid)

	mnu0 := app_db.Mnu_GetSubItemIdByType(1, "[D]")
	app_conf.SetVal("sub0_default", mnu0)
	mnu1 := app_db.Mnu_GetSubItemIdByType(2, "[D]")
	app_conf.SetVal("sub1_default", mnu1)
	mnu2 := app_db.Mnu_GetSubItemIdByType(3, "[D]")
	app_conf.SetVal("sub2_default", mnu2)
	mnu3 := app_db.Mnu_GetSubItemIdByType(4, "[D]")
	app_conf.SetVal("sub3_default", mnu3)
	mnu4 := app_db.Mnu_GetSubItemIdByType(5, "[D]")
	app_conf.SetVal("sub4_default", mnu4)
	mnu5 := app_db.Mnu_GetSubItemIdByType(6, "[D]")
	app_conf.SetVal("sub5_default", mnu5)

}