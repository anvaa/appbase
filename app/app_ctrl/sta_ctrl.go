package app_ctrl

import (
	"fmt"
	"strings"

	"app/app_db"
	"app/app_models"
	"srv/global"

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

func Sta_HistDelete(c *gin.Context) {

	staid := c.Param("id")

	if staid == "" {
		c.JSON(400, gin.H{"error": "missing id"})
		return
	}

	if err := staHist_Delete(staid); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "success"})

}

func sta_GetStatHistory(id any) ([]app_models.StatusHistory, error) {
	var stath []app_models.StatusHistory
	if err := app_db.AppDB.Preload("SubMenu0").Preload("Status").Preload("User").Where("item_id = ?", id).Order("updated_at desc").Find(&stath).Error; err != nil {
		return nil, err
	}
	return stath, nil
}

func Sta_GetStatNames() ([]app_models.Status, int) {
	var stats []app_models.Status
	app_db.AppDB.Order("staname").Find(&stats)
	return stats, len(stats)
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

func staHist_Delete(id any) error {
	// delete from status_history
	var stah app_models.StatusHistory
	if err := app_db.AppDB.Where("id = ?", id).Delete(&stah).Error; err != nil {
		return err
	}

	return nil

}

func Sta_HistAdd(c *gin.Context) {
	var body app_models.StatusHistory

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body.Note = strings.TrimSpace(body.Note)
	// fmt.Println("iem_id", body.ItemId, "user_id", body.UserId, "status_id", body.StatusId, body.Note)

	if err := sta_HistAdd(body); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "success", "url": "/search/" + global.IntToString(body.ItemId)})
}

func sta_HistAdd(newhist app_models.StatusHistory) error {
	// add new status history
	fmt.Println("iem_id", newhist.ItemId, "user_id", newhist.UserId, "status_id", newhist.StatusId, newhist.Note)
	if err := app_db.AppDB.Create(&newhist).Error; err != nil {
		return err
	}

	// update status on item
	// if err := app_db.AppDB.Model(&app_models.Items{}).Where("id = ?", newhist.ItemId).Updates(app_models.Items{
	// 	StatusId: newhist.StatusId,
	// }).Error; err != nil {
	// 	return err
	// }

	return nil
}
