package app_models

type Project struct {
	BaseModel
	Name string `gorm:"size:50" json:"name"`
	Note string `gorm:"size:255" json:"note"`

	StasubID int    `gorm:"default:0" json:"stasub_id"`
	Stasub   Stasub `gorm:"foreignKey:StasubID" json:"stasub"`

	TypsubID int    `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID;references:ID" json:"typsub"`

	//Persons   []Person  `gorm:"foreignKey:ProjectID" json:"persons"`
	//Phones    []Phone   `gorm:"foreignKey:ProjectID" json:"phones"`
	//Addresses []Address `gorm:"foreignKey:ProjectID" json:"addresses"`
	//Emails    []Email   `gorm:"foreignKey:ProjectID" json:"emails"`
}