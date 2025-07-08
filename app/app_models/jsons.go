package app_models

import "time"

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	UserUUID  int       `json:"user_uuid"`
	Locale    string    `json:"locale"`
	Expire    time.Time `json:"expire"`
}
type PrinterConf struct {
	Height   int `json:"height"`
	Width    int `json:"width"`
	Margin   int `json:"margin"`
	FontSize int `json:"font_size"`
	Prnttxt  int `json:"print_txt"`
}

type AppConf struct {
	StartPageFocus string `json:"start_page_focus"`
}

type PrintQR struct {
	UUID  string `json:"uuid"`
	Date  string `json:"date"`
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	Line3 string `json:"line3"`
}

type PrintBarCode struct {
	UUID  string `json:"uuid"`
	Line1 string `json:"line1"`
}
