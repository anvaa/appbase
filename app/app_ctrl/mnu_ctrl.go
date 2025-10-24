package app_ctrl

import (
	
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"app/app_db"
	"app/app_models"
)

func Mnu_UpdTitels(c *gin.Context) {
	var body struct {
		Mnu_id    int    `json:"mnu_id" binding:"required"`
		Sub_uuid  int    `json:"sub_uuid" binding:"required"`
		Mnu_title string `json:"mnu_title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println("Mnu_UpdTitels:", body)
	if body.Mnu_id != 500 || body.Sub_uuid == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID or submenu UUID"})
		return
	}

	body.Mnu_title = strings.TrimSpace(body.Mnu_title)

	// Update the menu title where uuid = body.Sub_uuid
	var menu app_models.Menu
	err := app_db.AppDB.Model(&menu).Where("uuid = ?", body.Sub_uuid).Update("title", body.Mnu_title).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}
