package app_api

import (
	"server/middleware"

	"app/app_ctrl"

	"github.com/gin-gonic/gin"
)

func tools_Api(r *gin.Engine) *gin.Engine {

	// Require authentication first
	r.Use(middleware.Verify)

	mnuGrp := r.Group("/menu", middleware.RequireLevel(30))
	{
		mnuGrp.POST("/addupd", app_ctrl.Sub_AddUpd)
		mnuGrp.POST("/delete", app_ctrl.Sub_Delete)
	}

	staGrp := r.Group("/status", middleware.RequireLevel(30))
	{
		staGrp.POST("/addupd", app_ctrl.Sta_AddUpd)
		staGrp.POST("/delete", app_ctrl.Sta_Delete)
	}

	typGrp := r.Group("/type", middleware.RequireLevel(30))
	{
		typGrp.POST("/addupd", app_ctrl.Typ_AddUpd)
		typGrp.POST("/delete", app_ctrl.Typ_Delete)
	}

	titelsGrp := r.Group("/title", middleware.RequireLevel(30))
	{
		titelsGrp.POST("/upd", app_ctrl.Mnu_UpdTitels)
	}

	toolsGrp := r.Group("/tools")
	{
		toolsGrp.GET("/titles", app_ctrl.ToolsTitles) // tools titles page
		toolsGrp.GET("/status", app_ctrl.ToolsStatus) // tools statuses page
		toolsGrp.GET("/types", app_ctrl.ToolsTypes)   // tools types page

	}

	return r
}
