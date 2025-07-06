package app_api

import (
	
	"srv/web/middleware"

	"app/app_ctrl"

	"github.com/gin-gonic/gin"
)

func tools_Api(r *gin.Engine) *gin.Engine {

	mnuGrp := r.Group("/menu")
	{
		mnuGrp.Use(middleware.RequireAuth)

		mnuGrp.POST("/addupd", app_ctrl.Sub_AddUpd)
		mnuGrp.POST("/delete", app_ctrl.Sub_Delete)

	}
	
	staGrp := r.Group("/status")
	{
		staGrp.Use(middleware.RequireAuth)

		staGrp.POST("/addupd", app_ctrl.Sta_AddUpd)
		staGrp.POST("/delete", app_ctrl.Sta_Delete)

	}

	typGrp := r.Group("/type")
	{
		typGrp.Use(middleware.RequireAuth)

		typGrp.POST("/addupd", app_ctrl.Typ_AddUpd)
		typGrp.POST("/delete", app_ctrl.Typ_Delete)

	}

	titelsGrp := r.Group("/title")
	{
		titelsGrp.Use(middleware.RequireAuth)

		titelsGrp.POST("/upd", app_ctrl.Mnu_UpdTitels)

	}

	toolsGrp := r.Group("/tools")
	{
		toolsGrp.Use(middleware.RequireAuth)

		toolsGrp.GET("/titles", app_ctrl.ToolsTitles)
		toolsGrp.GET("/conf", app_ctrl.ToolsConf) // tools configuration page

		toolsGrp.POST("/printconf", app_ctrl.PrintConf) // save print conf
		toolsGrp.POST("/appconf", app_ctrl.AppConf) //  application conifg/settings

	}

	return r
}