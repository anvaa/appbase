package app_models

type AppConf struct {
	StartPageFocus string `json:"start_page_focus"`
}

type Appinfo struct {
	Company     string
	AppName     string
	AppNameLong string
	Version     string
}

type Appbase struct {
	Title   string   `json:"title"`
	Logos   []string `json:"logos"`
	User    any      `json:"user"`
	Appinfo any      `json:"appinfo"`
}
