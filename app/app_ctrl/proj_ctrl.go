package app_ctrl

import (
	"app/app_db"
	"app/app_models"

	"strings"

	"fmt"

	"github.com/gin-gonic/gin"
)

func Proj_Edit(c *gin.Context) {

	project, err := app_db.GetProjectByUUID(c.Param("id"))
	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"title": "Project Not Found",
			"error": "The project you are trying to edit does not exist.",
		})
		return
	}

	c.HTML(200, "proj_edit.html", gin.H{
		"apibase": setApiBase(c),
		"js":      "projects.js",

		"project": project,
		"sta":     app_db.Sta_GetAllStasub(1),
		"typ":     app_db.Typ_GetAllTypsub(1),
	})
}

func Proj_Delete(c *gin.Context) {
	project, err := app_db.GetProjectByUUID(c.PostForm("id"))
	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"title": "Project Not Found",
			"error": "The project you are trying to delete does not exist.",
		})
		return
	}

	if err := app_db.AppDB.Delete(&project).Error; err != nil {
		c.HTML(500, "error.html", gin.H{
			"title": "Delete Error",
			"error": fmt.Sprintf("Failed to delete project: %v", err),
		})
		return
	}

	c.Redirect(302, "/app/start")
}

func Proj_AddUpd(c *gin.Context) {
	var body struct {
		UUID  int    `json:"uuid"`
		Name  string `json:"name"`
		Note  string `json:"note"`
		Staid int    `json:"staid"`
		Typid int    `json:"typid"`
		UID   int    `json:"uid"` // Current user ID
	}

	if err := c.BindJSON(&body); err != nil {
		c.HTML(400, "error.html", gin.H{
			"title": "Bad Request",
			"error": fmt.Sprintf("Invalid input: %v", err),
		})
		return
	}

	// Validate project name
	if body.Name == "" || body.Staid == 0 || body.Typid == 0 {
		c.HTML(400, "error.html", gin.H{
			"title": "Validation Error",
			"error": "Project cannot be empty.",
		})
		return
	}

	// Prepare project data
	_proj := app_models.Project{
		Name:     strings.TrimSpace(body.Name),
		Note:     strings.TrimSpace(body.Note),
		StasubID: body.Staid,
		TypsubID: body.Typid,
	}

	_proj.UpdatedBy = CurUserUUID // Set the current user ID

	if body.UUID > 0 {
		// Update existing project
		if err := app_db.AppDB.Model(&app_models.Project{}).Where("uuid = ?", body.UUID).Updates(_proj).Error; err != nil {
			c.HTML(500, "error.html", gin.H{
				"title": "Update Error",
				"error": fmt.Sprintf("Failed to update project: %v", err),
			})
			return
		}
	} else {
		// Add new project
		_proj.CreatedBy = CurUserUUID // Set the current user ID for creation
		
		if err := app_db.AppDB.Create(&_proj).Error; err != nil {
			c.HTML(500, "error.html", gin.H{
				"title": "Creation Error",
				"error": fmt.Sprintf("Failed to create project: %v", err),
			})
			return
		}
		body.UUID = _proj.UUID // Get the UUID of the newly created project
	}

	url := fmt.Sprintf("/proj/%d", body.UUID)
	c.JSON(200, gin.H{
		"message":  "success",
		"redirect": url,
	})
}
