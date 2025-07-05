package app_models

import (
	"time"
	
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type BaseModel struct {
	gorm.Model
	ID   int `gorm:"primaryKey,autoIncrement" json:"id"`
	UUID int `gorm:"index,unique" json:"uuid"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CreatedBy int `json:"created_by,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
	DeletedBy int `json:"deleted_by,omitempty"`
}

// BeforeCreate sets a UUID before creating a record.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.UUID = int(uuid.New().ID())
	return
}

type Users struct {
	BaseModel
	Email    string `gorm:"unique, size:255" json:"email"`
	Password string `gorm:"not null, size:255" json:"password"`
	Role     string `gorm:"default:user, size:20" json:"role"`
	IsAuth   bool   `gorm:"default:false" json:"is_auth"`
	Note     string `gorm:"size:255" json:"note"`
}

type Status struct {
	BaseModel
	Title  string   `gorm:"size:50" json:"title"`
	Stasub []Stasub `gorm:"foreignKey:StatusID" json:"stasub"`
}

type Stasub struct {
	BaseModel
	Name     string `gorm:"size:50" json:"name"`
	Type     string `gorm:"size:5" json:"type"`
	StatusID int    `json:"status_id"`
}
type Type struct {
	BaseModel
	Name   string   `gorm:"size:50" json:"name"`
	Typsub []Typsub `gorm:"foreignKey:TypeID" json:"typsub"`
}

type Typsub struct {
	BaseModel
	Name   string `gorm:"size:50" json:"name"`
	TypeID int    `json:"type_id"`
}

type Menu struct {
	BaseModel
	Title   string    `gorm:"size:50" json:"title"`
	Type    string    `gorm:"size:5" json:"type"`
	Menusub []Menusub `gorm:"foreignKey:MenuID" json:"menu_sub"`
}

type Menusub struct {
	BaseModel
	Name   string `gorm:"size:50" json:"name"`
	Type   string `gorm:"size:5" json:"type"`
	MenuID int    `json:"menu_id"`
}