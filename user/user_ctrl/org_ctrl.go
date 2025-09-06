package user_ctrl

import (
	"app/app_db"
	"app/app_models"
	
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
	orgUUID := c.Param("uuid")
	org, err := app_db.Org_GetByUUID(orgUUID)
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

func Org_Members(c *gin.Context) {
	orgUUID := c.Param("uuid")
	org, err := app_db.Org_GetByUUID(orgUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	usr, err := app_db.Users_GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// add only usr to nonmember if not in org.Users
	var nonMembers []*app_models.Users
	for _, u := range usr {
		if !containsUser(org.Users, u) {
			nonMembers = append(nonMembers, u)
		}
	}

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
	// fmt.Println("add member org_id:", body.OrgID, "user_id:", body.UserID)
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

	// Illegal to remove userID 1 admin
	if body.UserID == 1 {
		c.JSON(403, gin.H{"error": "Can not remove admin user"})
		return
	}

	//fmt.Println("rem memberorg_id:", body.OrgID, "user_id:", body.UserID)
	if err := app_db.Org_RemoveMember(body.OrgID, body.UserID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

// containsUser checks if a user is in the list of users
func containsUser(users []*app_models.Users, user *app_models.Users) bool {
	for _, u := range users {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}
