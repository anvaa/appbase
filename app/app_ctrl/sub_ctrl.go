package app_ctrl

import (
	"fmt"
	"net/http"
	"strings"

	"app/app_db"
	"app/app_models"

	"github.com/gin-gonic/gin"
)

func Sub_AddUpd(c *gin.Context) {
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

	var sub app_models.Menusub
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
		fmt.Println("Updating existing title sub:", body.Sub_uuid, sub.Name)
	}

	if body.Sub_uuid == 0 {
		// Check if the menu_id and name combo already exist
		if err := app_db.AppDB.Where("name = ? AND menu_id = ?", body.Val, body.Mnu_id).First(&sub).Error; err == nil {
			c.JSON(400, gin.H{"error": "SubItem already exists"})
			return
		}

		sub = app_models.Menusub{
			Name:   body.Val,
			MenuID: body.Mnu_id,
		}

		if err := app_db.AppDB.Create(&sub).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			fmt.Println("Error creating new menusub:", err)
			return
		}
		fmt.Println("Creating new menusub name:", body.Val, "for MenuID:", body.Mnu_id)
	}

	c.JSON(200, gin.H{"message": "success"})
}

func Sub_Delete(c *gin.Context) {
	var body struct {
		Sub_uuid int `json:"sub_uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var sub app_models.Menusub
	if err := app_db.AppDB.Where("uuid = ?", body.Sub_uuid).First(&sub).Error; err != nil {
		c.JSON(404, gin.H{"error": "SubItem not found"})
		return
	}

	if sub.Type == "[D]" {
		c.JSON(400, gin.H{"error": "Cannot delete default subitem"})
		return
	}

	if err := app_db.AppDB.Where("uuid = ?", body.Sub_uuid).Delete(&sub).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
