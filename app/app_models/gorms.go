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
	Name string `gorm:"size:20" json:"name"`
	Type string `gorm:"size:5" json:"type"`
}

type StatusHistory struct {
	BaseModel
	ItemId int `gorm:"foreignKey:ItemId" json:"item_id"`

	UserId int `gorm:"foreignKey:UserId" json:"user_id"`
	User   Users

	Note string `gorm:"size:255" json:"note"`

	StatusId int `gorm:"foreignKey:StatusId" json:"status_id"`
	Status   Status
}

// Menu represents a menu with subMenu.
type Menu struct {
	BaseModel
	Title    string    `gorm:"size:50" json:"title"`
	Type     string    `gorm:"size:5" json:"type"`
	SubItems []SubMenu `gorm:"foreignKey:MenuId" json:"sub_items"`
}

// SubMenu represents a submenu linked to menu.
type SubMenu struct {
	BaseModel
	Name   string `gorm:"size:50" json:"name"`
	Type   string `gorm:"size:5" json:"type"`
	MenuId int    `gorm:"foreignKey:MenuId" json:"menu_id"`
}
