package app_models

type Appinfo struct {
	Company     string
	CompanyURL  string
	AppName     string
	AppNameLong string
	Version     string
}

type Appbase struct {
	Title   string   `json:"title"`
	DbInfo  string   `json:"dbinfo"`
	Doindex bool     `json:"doindex"`
	Logos   []string `json:"logos"`
	User    any      `json:"user"`
	Appinfo any      `json:"appinfo"`
}

type DbConfig struct {
	Type     string
	Path     string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}
