package app_ctrl

import (
	"app/app_conf"
	"app/app_db"
	"app/app_models"

	"strings"

	"fmt"

	"github.com/gin-gonic/gin"
)

func Proj_View(c *gin.Context) {
	
	project, err := app_db.GetProjectByUUID(c.Param("id"))
	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"title": "Project Not Found",
			"error": "The project you are trying to view does not exist.",
		})
		return
	}

	c.HTML(200, "proj_view.html", gin.H{
		"apibase": setApiBase(c),
		"css":     "app.css",
		"js":      "proj.js",

		"project": project,
	})
}

func Proj_Edit(c *gin.Context) {

	project, err := app_db.GetProjectByUUID(c.Param("id"))
	
	if err != nil {

		c.HTML(404, "error.html", gin.H{
			"title": "Project Not Found",
			"error": "The project you are trying to edit does not exist.",
		})
		return
	}

	sta, err := app_db.Sta_GetAllStasub("Project")
	if err != nil {
		c.HTML(500, "error.html", gin.H{
			"title": "Internal Server Error",
			"error": fmt.Sprintf("Failed to retrieve status: %v", err),
		})
		return
	}

	typ, err := app_db.Typ_GetAllTypsub("Project")
	if err != nil {
		c.HTML(500, "error.html", gin.H{
			"title": "Internal Server Error",
			"error": fmt.Sprintf("Failed to retrieve types: %v", err),
		})
		return
	}

	c.HTML(200, "proj_edit.html", gin.H{
		"apibase": setApiBase(c),
		"js":      "proj.js",

		"project": project,
		"sta":     sta,
		"typ":     typ,
	})
}

func Proj_Delete(c *gin.Context) {
	project, err := app_db.GetProjectByUUID(c.Param("id"))
	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"title": "Project Not Found",
			"error": "The project you are trying to delete does not exist.",
		})
		return
	}

	project.DeletedByID = CurUserID // Set the deleter as the current user

	if err := app_db.AppDB.Delete(&project).Error; err != nil {
		c.HTML(500, "error.html", gin.H{
			"title": "Delete Error",
			"error": fmt.Sprintf("Failed to delete project: %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":  "success",
		"redirect": "/app/start",
	})
}

func Proj_AddUpd(c *gin.Context) {
	var body struct {
		UUID  uint    `json:"uuid"`
		Name  string `json:"name"`
		Staid uint    `json:"staid"`
		Typid uint    `json:"typid"`
		UID   uint    `json:"uid"` // Current user ID
	}

	if err := c.BindJSON(&body); err != nil {
		c.HTML(400, "error.html", gin.H{
			"title": "Bad Request",
			"error": fmt.Sprintf("Invalid input: %v", err),
		})
		return
	}

	// Validate project name
	if body.Name == "" {
		c.HTML(400, "error.html", gin.H{
			"title": "Validation Error",
			"error": "Project cannot be empty.",
		})
		return
	}

	// set default status and type if not provided
	if body.Staid == 0 || body.Typid == 0 {
		body.Staid = app_conf.GetUInt("stadef1")
		body.Typid = app_conf.GetUInt("typdef1")
	}

	// Prepare project data
	_proj := app_models.Project{
		Name:     strings.TrimSpace(body.Name),
		StasubID: body.Staid,
		TypsubID: body.Typid,
		UpdatedByID: CurUserID,
	}

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
		_proj.CreatedByID = CurUserID
		_proj.UpdatedByID = CurUserID // Set the updater as well
		_proj.DeletedByID = 1 // set default to admin
		fmt.Println(_proj.Name, _proj.StasubID, _proj.TypsubID, _proj.CreatedByID, _proj.UpdatedByID)
		
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
