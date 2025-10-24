package user_ctrl

import (
	"app/app_db"
	"app/app_models"
	"strings"

	"github.com/gin-gonic/gin"
)

func Org_View(c *gin.Context) {
	orgs, err := app_db.Org_GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.HTML(200, "org_view.html", gin.H{"orgs": orgs})
}

func Org_New(c *gin.Context) {
	c.HTML(200, "org_new.html", gin.H{"js": "org.js"})
}

func Org_Edit(c *gin.Context) {
	orgUUID := c.Param("uuid")
	org, err := app_db.Org_GetByUUID(orgUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.HTML(200, "org_edit.html", gin.H{"js": "org.js", "org": org})
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

func Org_Members(c *gin.Context) {
	orgUUID := c.Param("uuid")
	org, err := app_db.Org_GetByUUID(orgUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	users, err := app_db.Users_GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	nonMembers := filterNonMembers(org.Users, users)
	c.HTML(200, "org_members.html", gin.H{
		"js":    "org.js",
		"org":   org,
		"users": nonMembers,
	})
}

func Org_AddMember(c *gin.Context) {
	var body struct {
		OrgID  int `json:"org_id"`
		UserID int `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := app_db.Org_AddMember(body.OrgID, body.UserID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func Org_RemoveMember(c *gin.Context) {
	var body struct {
		OrgID  int `json:"org_id"`
		UserID int `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if body.UserID == 1 {
		c.JSON(403, gin.H{"error": "Can not remove admin user"})
		return
	}
	if err := app_db.Org_RemoveMember(body.OrgID, body.UserID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func filterNonMembers(members, allUsers []*app_models.Users) []*app_models.Users {
	var nonMembers []*app_models.Users
	memberMap := make(map[int]struct{}, len(members))
	for _, m := range members {
		memberMap[int(m.ID)] = struct{}{}
	}
	for _, u := range allUsers {
		if _, exists := memberMap[int(u.ID)]; !exists {
			nonMembers = append(nonMembers, u)
		}
	}
	return nonMembers
}
