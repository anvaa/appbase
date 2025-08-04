package app_models

type Project struct {
	BaseModel
	Name string `gorm:"size:50" json:"name"`
	Note string `gorm:"size:255" json:"note"`

	StasubID int    `gorm:"default:0" json:"stasub_id"`
	Stasub   Stasub `gorm:"foreignKey:StasubID" json:"stasub"`

	TypsubID int    `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID" json:"typsub"`

	CreatedbyID int `gorm:"default:0" json:"createdby_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedbyID" json:"created_by"`
	UpdatedbyID int `gorm:"default:0" json:"updatedby_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedbyID" json:"updated_by"`
	DeletedbyID int `gorm:"default:0" json:"deletedby_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedbyID" json:"deleted_by"`

	//Persons   []Person  `gorm:"foreignKey:ProjectID" json:"persons"`
	//Phones    []Phone   `gorm:"foreignKey:ProjectID" json:"phones"`
	//Addresses []Address `gorm:"foreignKey:ProjectID" json:"addresses"`
	//Emails    []Email   `gorm:"foreignKey:ProjectID" json:"emails"`
}
