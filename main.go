package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/parthav/restraunt-management/database"
	"github.com/parthav/restraunt-management/middleware"
	"github.com/parthav/restraunt-management/routes"
	"go.mongodb.org/mongo-driver/mongo"
    //"go.mongodb.org/mongo-driver/mongo/options"


	
)

var foodcollection *mongo.Collection = database.OpenCollection(database.Client,"food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	router := gin.New()

	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.MenuRoutes(router)
	routes.InvoiceRoutes(router)
	routes.FoodRoutes(router)


	router.Run(":"+port)
}
