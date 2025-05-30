package app_ctrl

import (
	"fmt"
	"strings"

	"app/app_db"
	"app/app_models"

	"github.com/gin-gonic/gin"
)

func Stat_GetStatusSubItems() ([]app_models.Status, int) {
	var stats []app_models.Status
	app_db.AppDB.Order("name").Find(&stats)
	return stats, len(stats)
}

func Sta_GetLatestStat(itmid any) string {
	var stah app_models.StatusHistory
	app_db.AppDB.Where("itmid = ?", itmid).Last(&stah)
	fmt.Println("Sta_GetLatestStat", app_db.Sta_GetNameById(stah.StatusId))
	return app_db.Sta_GetNameById(stah.StatusId)
}

func Sta_AddUpd(c *gin.Context) {
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

	if err := sta_AddUpd(body.Sub_uuid, body.Txt); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "success"})
}

func sta_AddUpd(uuid any, val string) error {
	// add/update name where uuid = ? in status
	var stat app_models.Status
	res := app_db.AppDB.Where("uuid = ?", uuid).
		Attrs(app_models.Status{
			Name: val,
		}).
		FirstOrCreate(&stat)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// update if exists
		err := app_db.AppDB.Model(&stat).Updates(app_models.Status{
			Name: val,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
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
	var sta app_models.Status
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

	c.JSON(200, gin.H{"message": "success",})
}