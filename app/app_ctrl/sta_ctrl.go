package app_ctrl

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"app/app_db"
	"app/app_models"

	"github.com/gin-gonic/gin"
)

func Sta_GetStatuses() []app_models.Status {
	var stats []app_models.Status
	err := app_db.AppDB.Preload("Stasub").Find(&stats).Error
	if err != nil {
		log.Println("Error getting statuses:", err)
		return nil
	}

	return sortStatuses(stats)
}

func sortStatuses(menu []app_models.Status) []app_models.Status {
	// Sort each menu's submenu by name ascending
	for i := range menu {
		sort.Slice(menu[i].Stasub, func(i2, j int) bool {
			return menu[i].Stasub[i2].Name < menu[i].Stasub[j].Name
		})
	}
	return menu
}

func Sta_AddUpd(c *gin.Context) {
	var body struct {
		Mnu_id   int    `json:"mnu_id"`
		Sub_uuid int    `json:"sub_uuid"`
		Val      string `json:"val"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body.Val = strings.TrimSpace(body.Val)

	var sub app_models.Stasub
	if body.Sub_uuid != 0 {
		// If Sub_uuid exists, update
		if err := app_db.AppDB.Where("uuid = ?", body.Sub_uuid).First(&sub).Error; err != nil {
			c.JSON(404, gin.H{"error": "Stasub not found"})
			return
		}

		sub.Name = body.Val
		if err := app_db.AppDB.Save(&sub).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Updating existing status:", body.Sub_uuid, sub.Name)
	}

	if body.Sub_uuid == 0 {
		// Check if the menu_id and name combo already exist
		if err := app_db.AppDB.Where("name = ? AND status_id = ?", body.Val, body.Mnu_id).First(&sub).Error; err == nil {
			c.JSON(400, gin.H{"error": "SubItem already exists"})
			return
		}

		sub = app_models.Stasub{
			Name:     body.Val,
			StatusID: body.Mnu_id,
		}

		if err := app_db.AppDB.Create(&sub).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			fmt.Println("Error creating new stasub:", err)
			return
		}
		fmt.Println("Creating new stasub name:", body.Val, "for StatusID:", body.Mnu_id)
	}

	c.JSON(200, gin.H{"message": "success"})
}

func Sta_Delete(c *gin.Context) {
	var body struct {
		Sub_uuid int `json:"sub_uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if body.Sub_uuid == 0 {
		c.JSON(400, gin.H{"error": "missing UUID"})
		return
	}

	// if not type = "[D", delete
	var sta app_models.Stasub
	if err := app_db.AppDB.Where("uuid = ?", body.Sub_uuid).First(&sta).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if sta.Type == "[D]" {
		c.JSON(400, gin.H{"error": "cannot delete default status"})
		return
	}

	// delete from status
	if err := app_db.AppDB.Where("uuid = ?", body.Sub_uuid).Delete(&sta).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}
