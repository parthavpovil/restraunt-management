package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parthav/restraunt-management/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orders",controllers.GetOrders())
	incomingRoutes.GET("/orders/:order_id",controllers.GetOrder())
	incomingRoutes.POST("/order",controllers.CreateOrder())
	incomingRoutes.PATCH("/order/:order_id",controllers.UpdateOrder())
}