package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parthav/restraunt-management/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/menus",controllers.GetMenus())
	incomingRoutes.GET("/menus/:menu_id",controllers.GetMenu())
	incomingRoutes.PUT("/menus",controllers.CreateMenu())
	incomingRoutes.PATCH("/menus/:id",controllers.UpdateMenu())

}