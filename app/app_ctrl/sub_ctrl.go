package app_ctrl

import (
	"net/http"
	"strings"

	"app/app_db"
	"app/app_models"

	"github.com/gin-gonic/gin"
)

func Sub_AddUpd(c *gin.Context) {
	var body struct {
		Mnu_id   int    `json:"mnu_id" binding:"required"`
		Sub_uuid int    `json:"sub_uuid"`
		Txt      string `json:"txt"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body.Txt = strings.TrimSpace(body.Txt)


	// If Sub_uuid not 0, means we are updating an existing item: 
		// return after update with status ok
	
	// If Sub_uuid is 0, means we are creating a new item:
		// Check if the menuID and Name combo already exist
		// If found, update the existing item
		// If not found, create new item

	// Check if the menu ID and name combo already exist
	var sub app_models.SubMenu
	if body.Sub_uuid != 0 {
		if err := app_db.AppDB.Where("uuid = ?", body.Sub_uuid).First(&sub).Error; err != nil {
			c.JSON(404, gin.H{"error": "SubItem not found"})
			return
		}

		sub.Name = body.Txt
		if err := app_db.AppDB.Save(&sub).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} 

	if body.Sub_uuid == 0 {
		// Check if the menu_id and name combo already exist
		if err := app_db.AppDB.Where("menu_id = ? AND name like ?", body.Mnu_id, body.Txt).First(&sub).Error; err == nil {
			// If found, update the existing item
			sub.Name = body.Txt
			if err := app_db.AppDB.Save(&sub).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		} else {
			// If not found, create a new item
			sub = app_models.SubMenu{
				MenuId: body.Mnu_id,
				Name:   body.Txt,
			}
			if err := app_db.AppDB.Create(&sub).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}


	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func Sub_Delete(c *gin.Context) {
	var body struct {
		Sub_uuid int `json:"sub_uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var sub app_models.SubMenu
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
