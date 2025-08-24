package app_models

import "time"

type Project struct {
	BaseModel
	Name string `gorm:"size:50" json:"name"`

	StasubID int    `gorm:"default:0" json:"stasub_id"`
	Stasub   Stasub `gorm:"foreignKey:StasubID" json:"stasub"`

	TypsubID int    `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID" json:"typsub"`

	CreatedbyID int   `gorm:"default:1" json:"createdby_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedbyID" json:"created_by"`
	UpdatedbyID int   `gorm:"default:1" json:"updatedby_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedbyID" json:"updated_by"`
	DeletedbyID int   `gorm:"default:1" json:"deletedby_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedbyID" json:"deleted_by"`

	Notes []Notes `gorm:"foreignKey:ProjectID" json:"notes"`
	//Persons   []Person  `gorm:"foreignKey:ProjectID" json:"persons"`
	//Phones    []Phone   `gorm:"foreignKey:ProjectID" json:"phones"`
	//Addresses []Address `gorm:"foreignKey:ProjectID" json:"addresses"`
	//Emails    []Email   `gorm:"foreignKey:ProjectID" json:"emails"`
}

type Notes struct {
	BaseModel
	Header  string `gorm:"size:50, default:'Topic'" json:"header"`
	Content string `gorm:"type:text" json:"content"`

	ProjectID int `gorm:"default:0" json:"project_id"`
	PersonID  int `gorm:"default:0" json:"person_id"`

	TypsubID int    `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID" json:"typsub"`

	CreatedbyID int   `gorm:"default:1" json:"createdby_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedbyID" json:"created_by"`
	UpdatedbyID int   `gorm:"default:1" json:"updatedby_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedbyID" json:"updated_by"`
	DeletedbyID int   `gorm:"default:1" json:"deletedby_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedbyID" json:"deleted_by"`
}

type Person struct {
	BaseModel
	FirstName   string    `gorm:"size:100" json:"first_name"`
	LastName    string    `gorm:"size:100" json:"last_name"`
	DOB         time.Time `json:"dob"`
	Gender      string    `gorm:"size:10" json:"gender"`
	Nationality string    `gorm:"size:50" json:"nationality"`

	Notes []Notes `gorm:"foreignKey:PersonID" json:"notes"`

	ProjectID int `gorm:"default:0" json:"project_id"`

	StasubID int    `gorm:"default:0" json:"stasub_id"`
	Stasub   Stasub `gorm:"foreignKey:StasubID" json:"stasub"`
	TypsubID int    `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID" json:"typsub"`

	CreatedbyID int   `gorm:"default:1" json:"createdby_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedbyID" json:"created_by"`
	UpdatedbyID int   `gorm:"default:1" json:"updatedby_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedbyID" json:"updated_by"`
	DeletedbyID int   `gorm:"default:1" json:"deletedby_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedbyID" json:"deleted_by"`
}

type Email struct {
	BaseModel
	Address  string `gorm:"size:255;unique" json:"address"`
	Password string `gorm:"size:255" json:"password"`

	ProjectID int `gorm:"default:0" json:"project_id"`

	CreatedbyID int   `gorm:"default:1" json:"createdby_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedbyID" json:"created_by"`
	UpdatedbyID int   `gorm:"default:1" json:"updatedby_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedbyID" json:"updated_by"`
	DeletedbyID int   `gorm:"default:1" json:"deletedby_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedbyID" json:"deleted_by"`
}
