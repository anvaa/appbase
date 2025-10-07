package app_models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      uint           `gorm:"index;unique" json:"uuid"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy uint           `gorm:"default:0" json:"created_by,omitempty"`
	UpdatedBy uint           `gorm:"default:0" json:"updated_by,omitempty"`
	DeletedBy uint           `gorm:"default:0" json:"deleted_by,omitempty"`
}

// SetAuditFields sets audit fields based on the action and user ID.
func (b *BaseModel) SetAuditFields(action string, userID uint) {
	switch action {
	case "create":
		b.CreatedBy = userID
	case "update":
		b.UpdatedBy = userID
	case "delete":
		b.DeletedBy = userID
	}
}

// BeforeCreate sets CreatedBy and UUID before creating a record.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if userID, ok := tx.Statement.Context.Value("user_id").(uint); ok {
		b.CreatedBy = userID
	}
	b.UUID = uint(uuid.New().ID())
	return
}

// BeforeUpdate sets UpdatedBy before updating a record.
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	if userID, ok := tx.Statement.Context.Value("user_id").(uint); ok {
		b.UpdatedBy = userID
	}
	return
}

// BeforeDelete sets DeletedBy before deleting a record.
func (b *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
	if userID, ok := tx.Statement.Context.Value("user_id").(uint); ok {
		b.DeletedBy = userID
	}
	return
}

type Users struct {
	BaseModel
	Username     string    `gorm:"unique;size:255" json:"username"`
	Email        string    `gorm:"unique;size:255" json:"email"`
	Password     string    `gorm:"not null;size:255" json:"password"`
	IsAuth       bool      `gorm:"default:false" json:"is_auth"`
	LastLogin    time.Time `gorm:"autoUpdateTime" json:"last_login"`
	Note         string    `gorm:"type:text" json:"note"`
	TokenVersion int       `gorm:"default:1" json:"token_version"`

	AuthLevelID int       `gorm:"default:5" json:"auth_level_id"`
	AuthLevel   AuthLevel `gorm:"foreignKey:AuthLevelID" json:"auth_level"`

	Org []*Org `gorm:"many2many:user_orgs;" json:"orgs,omitempty"`
}

type Org struct {
	BaseModel
	Name  string   `gorm:"unique;size:255" json:"name"`
	Note  string   `gorm:"type:text" json:"note"`
	Users []*Users `gorm:"many2many:user_orgs;" json:"users"`
}

type AuthLevel struct {
	BaseModel
	Name  string `gorm:"size:50" json:"name"`
	Level int    `gorm:"default:10" json:"level"`
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
