package app_conf

import (
	"time"
	"fmt"
)

var StartTime int64

var Company string = "Raadig AS"
var AppName string = "Alias"
var AppNameLong string = "- Ditt andre liv -"
var Version string = time.Now().Format("060102")

type appInfo struct {
	Company     string
	AppName     string
	AppNameLong string
	Version     string
}

func AppInfo() appInfo {
	return appInfo{Company, AppName, AppNameLong, Version}
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