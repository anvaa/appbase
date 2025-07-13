package app_ctrl

import (
	"app/app_db"
	"app/app_models"
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
 
	fmt.Println("Editing project:", project.UUID, project.Name)
	c.HTML(200, "proj_edit.html", gin.H{
		"title":   appname + " - Edit Project",
		"appinfo": appinfo,
		"js":      "projects.js",
		"user":    c.Keys["user"],
		
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

func Proj_Update (c *gin.Context) {
	project := app_models.Project{}

	if project.Name == "" {
		c.HTML(400, "error.html", gin.H{
			"title": "Validation Error",
			"error": "Project name cannot be empty.",
		})
		return
	}

	if err := app_db.AppDB.Save(&project).Error; err != nil {
		c.HTML(500, "error.html", gin.H{
			"title": "Database Error",
			"error": fmt.Sprintf("Failed to save project: %v", err),
		})
		return
	}

	c.Redirect(302, "/app/start")
}