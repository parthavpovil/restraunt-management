package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parthav/restraunt-management/controllers"
)

func InvoiceRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/invoices",controllers.GetInvoices())
	incomingRoutes.GET("/invoices/:invoice_id",controllers.GetInvoice())
	incomingRoutes.POST("/invoices",controllers.CreateInvoice())
	incomingRoutes.PATCH("/invoices/:id",controllers.UpdateInvoice())
}