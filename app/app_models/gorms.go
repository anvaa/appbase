package app_models

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type BaseModel struct {
	gorm.Model
	ID   int `gorm:"primaryKey,autoIncrement" json:"id"`
	UUID int `gorm:"index,unique" json:"uuid"`
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
	StaSub []StaSub `gorm:"foreignKey:StatusID" json:"sta_sub"`
}

type StaSub struct {
	BaseModel
	Name     string `gorm:"size:50" json:"name"`
	Type     string `gorm:"size:5" json:"type"`
	StatusID int    `gorm:"foreignKey:StatusID" json:"status_id"`
}

// Menu represents a menu with subMenu.
type Menu struct {
	BaseModel
	Title   string    `gorm:"size:50" json:"title"`
	Type    string    `gorm:"size:5" json:"type"`
	MenuSub []MenuSub `gorm:"foreignKey:MenuID" json:"menu_sub"`
}

// SubMenu represents a submenu linked to menu.
type MenuSub struct {
	BaseModel
	Name   string `gorm:"size:50" json:"name"`
	Type   string `gorm:"size:5" json:"type"`
	MenuID int    `gorm:"foreignKey:MenuID" json:"menu_id"`
}
