package app_models

import "time"

// BaseModel should be defined elsewhere, assumed to include ID, CreatedAt, UpdatedAt, etc.

type Project struct {
	BaseModel
	Name string `gorm:"size:50" json:"name"`

	StasubID uint   `gorm:"default:0" json:"stasub_id"`
	Stasub   Stasub `gorm:"foreignKey:StasubID" json:"stasub"`

	TypsubID uint   `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID" json:"typsub"`

	CreatedByID uint  `gorm:"default:1" json:"created_by_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"created_by"`
	UpdatedByID uint  `gorm:"default:1" json:"updated_by_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"updated_by"`
	DeletedByID uint  `gorm:"default:1" json:"deleted_by_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"deleted_by"`

	Notes   []Notes  `gorm:"foreignKey:ProjectID" json:"notes"`
	Persons []Person `gorm:"foreignKey:ProjectID" json:"persons"`
	// Phones    []Phone   `gorm:"foreignKey:ProjectID" json:"phones"`
	// Addresses []Address `gorm:"foreignKey:ProjectID" json:"addresses"`
	Emails []Emails `gorm:"foreignKey:ProjectID" json:"emails"`
}

type Person struct {
	BaseModel
	FirstName   string    `gorm:"size:100" json:"first_name"`
	LastName    string    `gorm:"size:100" json:"last_name"`
	DOB         time.Time `json:"dob"`
	Gender      string    `gorm:"size:10" json:"gender"`
	Nationality string    `gorm:"size:50" json:"nationality"`

	Notes  []Notes  `gorm:"foreignKey:PersonID" json:"notes"`
	Emails []Emails `gorm:"foreignKey:PersonID" json:"emails"`

	ProjectID uint `gorm:"default:0" json:"project_id"`

	StasubID uint   `gorm:"default:0" json:"stasub_id"`
	Stasub   Stasub `gorm:"foreignKey:StasubID" json:"stasub"`
	TypsubID uint   `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID" json:"typsub"`

	CreatedByID uint  `gorm:"default:1" json:"created_by_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"created_by"`
	UpdatedByID uint  `gorm:"default:1" json:"updated_by_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"updated_by"`
	DeletedByID uint  `gorm:"default:1" json:"deleted_by_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"deleted_by"`
}

type Notes struct {
	BaseModel
	Header  string `gorm:"size:50;" json:"header"`
	Content string `gorm:"type:text" json:"content"`

	ProjectID uint `gorm:"default:0" json:"project_id"`
	PersonID  uint `gorm:"default:0" json:"person_id"`

	TypsubID uint   `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID" json:"typsub"`

	CreatedByID uint  `gorm:"default:1" json:"created_by_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"created_by"`
	UpdatedByID uint  `gorm:"default:1" json:"updated_by_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"updated_by"`
	DeletedByID uint  `gorm:"default:1" json:"deleted_by_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"deleted_by"`
}

type Emails struct {
	BaseModel
	Email       string `gorm:"size:255;unique" json:"email"`
	Password    string `gorm:"size:255" json:"password"`
	Nationality string `gorm:"size:50" json:"nationality"`

	ProjectID uint `gorm:"default:0" json:"project_id"`
	PersonID  uint `gorm:"default:0" json:"person_id"`

	StasubID uint   `gorm:"default:0" json:"stasub_id"`
	Stasub   Stasub `gorm:"foreignKey:StasubID" json:"stasub"`
	TypsubID uint   `gorm:"default:0" json:"typsub_id"`
	Typsub   Typsub `gorm:"foreignKey:TypsubID" json:"typsub"`

	CreatedByID uint  `gorm:"default:1" json:"created_by_id"`
	CreatedBy   Users `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"created_by"`
	UpdatedByID uint  `gorm:"default:1" json:"updated_by_id"`
	UpdatedBy   Users `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"updated_by"`
	DeletedByID uint  `gorm:"default:1" json:"deleted_by_id"`
	DeletedBy   Users `gorm:"foreignKey:DeletedByID;constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT" json:"deleted_by"`
}
