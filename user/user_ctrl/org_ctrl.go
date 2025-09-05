package user_ctrl

import (
	"app/app_db"

	"strings"

	"github.com/gin-gonic/gin"
)

func Org_View(c *gin.Context) {
	// Get all organizations
	orgs, err := app_db.Org_GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "org_view.html", gin.H{
		"orgs": orgs,
	})
}

func Org_New(c *gin.Context) {
	c.HTML(200, "org_new.html", gin.H{
		"js": "org.js",
	})
}

func Org_Edit(c *gin.Context) {
	orgID := c.Param("uuid")
	org, err := app_db.Org_GetByUUID(orgID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "org_edit.html", gin.H{
		"js":  "org.js",
		"org": org,
	})
}

func Org_AddUpd(c *gin.Context) {
	var body struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
		Note string `json:"note"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body.Name = strings.TrimSpace(body.Name)
	body.Note = strings.TrimSpace(body.Note)
	
	if err := app_db.Org_AddUpd(body.UUID, body.Name, body.Note); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func Org_Delete(c *gin.Context) {
	orgID := c.Param("uuid")
	if err := app_db.Org_Delete(orgID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}
