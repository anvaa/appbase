package app_models
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

type TableActions struct {
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	DeletedBy string `json:"deleted_by"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	Deleted   string `json:"deleted"`
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
