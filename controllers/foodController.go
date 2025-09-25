package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/parthav/restraunt-management/database"
	"github.com/parthav/restraunt-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		foodId := c.Param("food_id")

		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error  while fetching food",
			})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var food models.Food
		var menu models.Menu

		err := c.BindJSON(&food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return 
		}
		validationError :=validate.Struct(food)
		if validationError !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":validationError.Error(),
			})
			return 
		}
		err =menuCollection.FindOne(ctx,bson.M{"menu-id":food.Menu_id}).Decode(&menu)
		if err !=nil{
			msg:= fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":msg,
			})
			return 
		}
		food.Created_at=time.Now()
		food.Updated_at=time.Now()
		food.ID =primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = ToFixed(*food.Price,2)
		food.Price = &num

		result, err :=foodCollection.InsertOne(ctx,food)
		if err !=nil{
			msg :=fmt.Sprintf("fodd item was not added")
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":msg,
			})
			return 
		}
		c.JSON(http.StatusOK,result)



	}
}

func round(num float64) int {

}

func ToFixed(num float64, precision int) float64 {

}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
