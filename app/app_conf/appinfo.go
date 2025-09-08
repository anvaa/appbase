package app_conf

import (
	"fmt"
	"time"

	"app/app_models"
)

var (
	StartTime int64
	version   string
)

func AppInfo() app_models.Appinfo {

	version = GetString("app_version")
	if version == "" {
		version = time.Now().Format("060102")
	}

	return app_models.Appinfo{
		Company:     GetString("company_name"),
		CompanyURL:  GetString("company_url"),
		AppName:     GetString("app_name"),
		AppNameLong: GetString("app_name_long"),
		Version:     version,
	}
}

func RunTime() string {

	runt := time.Unix(StartTime, 0)
	timediff := time.Since(runt)
	sincet := runt.Format("2006-01-02 15:04")

	// timediff in years, dys, hours and minutes
	years := timediff / (365 * 24 * time.Hour)
	timediff -= years * 365 * 24 * time.Hour
	days := timediff / (24 * time.Hour)
	timediff -= days * 24 * time.Hour
	hours := timediff / time.Hour
	timediff -= hours * time.Hour
	minutes := timediff / time.Minute

	timetxt := fmt.Sprintf("Running for %d years, %d days, %d hours, %d minutes since %s", years, days, hours, minutes, sincet)

	// if less than 1 minute
	if minutes == 0 {
		timetxt = fmt.Sprintf("Started just now @%s", runt.Format("2006-01-02 15:04:15"))
	}

	return timetxt
}

func GetLocalTime() string {

	// get the server current timezone
	timezone, _ := time.Now().Local().Zone()

	// get the local time
	loctime := time.Now().Format("2006-01-02 15:04:05")

	timetxt := fmt.Sprintf("%s %s", loctime, timezone)
	return timetxt
}
