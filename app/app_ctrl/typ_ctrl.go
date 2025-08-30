package app_ctrl

import (
	"app/app_db"
	"app/app_models"

	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func Typ_GetTypes() []app_models.Type {
	var types []app_models.Type
	err := app_db.AppDB.Preload("Typsub").Find(&types).Error
	if err != nil {
		return nil
	}
	// Sort each type's subtypes by name ascen		ng
	for i := range types {
		sort.Slice(types[i].Typsub, func(i2, j int) bool {
			return types[i].Typsub[i2].Name < types[i].Typsub[j].Name
		})
	}
	return types
}

func Typ_Delete(c *gin.Context) {
	var body struct {
		Sub_uuid int `json:"sub_uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println("Deleting type with UUID:", body.Sub_uuid)
	if body.Sub_uuid == 0 {
		c.JSON(400, gin.H{"error": "missing UUID"})
		return
	}

	// delete from typesub
	var typ app_models.Typsub
	if err := app_db.AppDB.Where("uuid = ?", body.Sub_uuid).Delete(&typ).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}

func Typ_AddUpd(c *gin.Context) {
	var body struct {
		Typ_id   int    `json:"mnu_id"`
		Sub_uuid int    `json:"sub_uuid"`
		Val      string `json:"val"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body.Val = strings.TrimSpace(body.Val)
	fmt.Println("Adding type", body.Val, body.Sub_uuid, body.Typ_id)
	var sub app_models.Typsub
	if body.Sub_uuid != 0 {
		// If Sub_uuid exists, update
		if err := app_db.AppDB.Where("uuid = ?", body.Sub_uuid).First(&sub).Error; err != nil {
			c.JSON(404, gin.H{"error": "Typesub not found"})
			return
		}

		sub.Name = body.Val
		if err := app_db.AppDB.Save(&sub).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Updating existing type:", body.Sub_uuid, sub.Name)
	}

	if body.Sub_uuid == 0 {
		// Check if the menu_id and name combo already exist
		if err := app_db.AppDB.Where("name = ? AND type_id = ?", body.Val, body.Typ_id).First(&sub).Error; err == nil {
			c.JSON(400, gin.H{"error": "SubItem already exists"})
			return
		}

		sub = app_models.Typsub{
			Name:   body.Val,
			TypeID: body.Typ_id,
		}

		if err := app_db.AppDB.Create(&sub).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			fmt.Println("Error creating new typesub:", err)
			return
		}
		fmt.Println("Creating new typesub name:", body.Val, "for TypeID:", body.Typ_id)
	}

	c.JSON(200, gin.H{"message": "success"})
}
