package app_ctrl

import (
	"app/app_conf"
	"app/app_models"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

func AppConf(c *gin.Context) {
	var body app_models.AppConf

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(body.StartPageFocus)

	app_conf.SetVal("start_page_focus", body.StartPageFocus)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "App configuration updated successfully",
		"data":    body,
	})
}