package user_ctrl

import (
	"app/app_models"
	"net/http"
	"server/srv_conf"

	"github.com/gin-gonic/gin"
)

func DB_SaveDbConf(c *gin.Context) {
	var newConfig app_models.DbConfig

	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := srv_conf.SetDBConfig(newConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"msg":     "Database configuration saved successfully.",
	})
}
