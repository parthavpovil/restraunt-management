package controllers

import (
	"context"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parthav/restraunt-management/database"
	"github.com/parthav/restraunt-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var orderCollection *mongo.Collection =database.OpenCollection(database.Client,"order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		 ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		 
		result, err :=orderCollection.Find(ctx,bson.M{})

		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error occured while lisitng order item",
			})
			return 
		}
		var allOrders []bson.M
		err = result.All(ctx,&allOrders)

		if err !=nil{
			c.JSON(http.StatusInternalServerError, gin.H{
                "error": "error occurred while getting order items",
            })
            return
		}
		c.JSON(http.StatusOK,allOrders)
		
}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var order models.Order
		orderID := c.Param("order_id")

		err :=orderCollection.FindOne(ctx,bson.M{"order_id": orderID}).Decode(&order)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error while fetching order",
			})
			return 
		}
		c.JSON(http.StatusOK,order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx , cancel :=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var table models.Table
		var order models.Order

		err := c.BindJSON(&order)

		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}
		validationErr :=validate.Struct(order)

		if validationErr !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":validationErr.Error(),
			})
			return 
		}

		if order.Table_id !=nil{
			tableCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)
			if err !=nil{
				c.JSON(http.StatusInternalServerError,gin.H{
					"error":"table was not found",
				})
				return
			}
			order.Created_at=time.Now()
			order.Updated_at=time.Now()

			order.ID=primitive.NewObjectID()
			order.Order_id=order.ID.Hex()
			result,err :=orderCollection.InsertOne(ctx,order)
			if err !=nil{
				c.JSON(http.StatusInternalServerError,gin.H{
					"Error":"order was not created",
				})
				return 
			}
			c.SecureJSON(http.StatusOK,result)
		}


	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var table models.Table
		var order models.Order

		var updateObj primitive.D

		orderID :=c.Param("order_id")

		err :=c.BindJSON(&order)

		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}
		if order.Table_id !=nil{
			err :=menuCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)

			if err!=nil{
				c.JSON(http.StatusInternalServerError,gin.H{
					"error":"error while fetching table information tbale not found",
				})
				return 
			}
			updateObj=append(updateObj, bson.E{"table_id",order.Table_id})
		}
		order.Updated_at=time.Now()
		updateObj=append(updateObj, bson.E{"updated_at",order.Updated_at})

		upsert :=true

		filter := bson.M{"order_id":orderID}

		opt := options.UpdateOptions{
			Upsert: &upsert,

		}
		result,err :=orderCollection.UpdateOne(ctx,
			filter,
		bson.D{
			{"$set",updateObj},
		},
		&opt,)

		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error while updating the order",
			})
			return 
		}
		c.JSON(http.StatusOK,result)

	}
}

func OrderItemOrderCreator(order models.Order) string{
	order.Created_at= time.Now()
	order.Updated_at =time.Now()
	order.ID =primitive.NewObjectID()
	order.Order_id=order.ID.Hex()
	ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

	orderCollection.InsertOne(ctx,order)
	return order.Order_id
}