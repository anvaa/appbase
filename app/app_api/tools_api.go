package app_api

import (
	"app/app_ctrl"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func ToolsApi(r *gin.Engine) *gin.Engine {
	// Menu routes - require level 30
	menu := r.Group("/menu")
	menu.Use(middleware.Verify, middleware.RequireLevel(30))
	menu.POST("/addupd", app_ctrl.Sub_AddUpd)
	menu.POST("/delete", app_ctrl.Sub_Delete)

	// Status routes - require level 30
	status := r.Group("/status")
	status.Use(middleware.Verify, middleware.RequireLevel(30))
	status.POST("/addupd", app_ctrl.Sta_AddUpd)
	status.POST("/delete", app_ctrl.Sta_Delete)

	// Type routes - require level 30
	typ := r.Group("/type")
	typ.Use(middleware.Verify, middleware.RequireLevel(30))
	typ.POST("/addupd", app_ctrl.Typ_AddUpd)
	typ.POST("/delete", app_ctrl.Typ_Delete)

	// Title routes - require level 30
	title := r.Group("/title")
	title.Use(middleware.Verify, middleware.RequireLevel(30))
	title.POST("/upd", app_ctrl.Mnu_UpdTitels)

	// Tools routes - require level 30
	tools := r.Group("/tools")
	tools.Use(middleware.Verify, middleware.RequireLevel(30))
	tools.GET("/titles", app_ctrl.ToolsTitles)
	tools.GET("/status", app_ctrl.ToolsStatus)
	tools.GET("/types", app_ctrl.ToolsTypes)

	return r
}
