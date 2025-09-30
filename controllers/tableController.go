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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client,"tables")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var allTables []bson.M

		result,err :=tableCollection.Find(ctx,bson.M{})
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error retrivng tables",
			})
			return 
		}
		err =result.All(ctx,&allTables)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error retrivng tables",
			})
			return 
		}
		c.JSON(http.StatusOK,allTables)


	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel :=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		tableId :=c.Param("table_id")
		var table models.Table

		err :=tableCollection.FindOne(ctx,bson.M{"table_id":tableId}).Decode(&table)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error retriving table from db",
			})
			return 
		}
		c.JSON(http.StatusOK,table)


	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var table models.Table
		err :=c.BindJSON(&table)
		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}

		validErr := validate.Struct(table)
		if validErr !=nil{
				c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}
		table.ID=primitive.NewObjectID()
		table.Created_at=time.Now()
		table.Table_id=table.ID.Hex()
		table.Updated_at=time.Now()

		result,err :=tableCollection.InsertOne(ctx,&table)
		if err !=nil{
				c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error inserting table to db",
			})
			return 
		}
		c.JSON(http.StatusOK,result)


		
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var table models.Table
		var updateObj primitive.D

		var tableId =c.Param("table_id")

		err:=c.BindJSON(&table)
		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}	
		if table.Number_of_guests !=nil{
			updateObj=append(updateObj, bson.E{"number_of_guests",table.Number_of_guests})
		}

		if table.Table_number !=nil{
			updateObj=append(updateObj, bson.E{"table_number",table.Table_number})
		}

		table.Updated_at=time.Now()

		upsert :=true
		opt :=options.UpdateOptions{
			Upsert: &upsert,

		}
		filter :=bson.M{"table_id":tableId}

		result,err :=tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)
		if err!=nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"erro updating table",
			})
			return
		}
		c.JSON(http.StatusOK,result)

	}
}
