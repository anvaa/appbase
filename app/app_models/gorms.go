package app_models

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type BaseModel struct {
	gorm.Model
	ID   uint `gorm:"primaryKey,autoIncrement" json:"id"`
	UUID uint `gorm:"index,unique" json:"uuid"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	CreatedBy uint `gorm:"default:0" json:"created_by,omitempty"`
	UpdatedBy uint `gorm:"default:0" json:"updated_by,omitempty"`
	DeletedBy uint `gorm:"default:0" json:"deleted_by,omitempty"`
}

// BeforeCreate sets a UUID before creating a record.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.UUID = uint(uuid.New().ID())
	return
}

type Users struct {
	BaseModel
	Email        string    `gorm:"unique, size:255" json:"email"`
	Password     string    `gorm:"not null, size:255" json:"password"`
	Orgname      string    `gorm:"size:255" json:"orgname"`
	IsAuth       bool      `gorm:"default:false" json:"is_auth"`
	LastLogin    time.Time `gorm:"autoUpdateTime" json:"last_login"`
	Note         string    `gorm:"size:255" json:"note"`
	TokenVersion int       `gorm:"default:1" json:"token_version"`

	AuthLevelID int       `gorm:"default:3" json:"auth_level_id"` // user
	AuthLevel   AuthLevel `gorm:"foreignKey:AuthLevelID" json:"auth_level"`
}

type AuthLevel struct {
	BaseModel
	Name  string `gorm:"size:50" json:"name"`
	Level int    `gorm:"default:3" json:"level"` // 1=user, 5=admin, 10=superuser
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
	Title  string   `gorm:"size:50" json:"title"`
	Typsub []Typsub `gorm:"foreignKey:TypeID" json:"typsub"`
}

type Typsub struct {
	BaseModel
	Name   string `gorm:"size:50" json:"name"`
	Type   string `gorm:"size:5" json:"type"`
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
