package controllers

import (
	"context"
	"fmt"
	"log"
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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while lisitng menu items",
			})
			return
		}
		var allMenu []bson.M
		if err = result.All(ctx, &allMenu); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenu)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		menuId := c.Param("menu_id")

		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erorr": "error fetching menu",
			})
			return
		}

		c.JSON(http.StatusOK, menu)

	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

		var menu models.Menu
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := c.BindJSON(&menu)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return

		}
		menu.Created_at = time.Now()
		menu.Updated_at = time.Now()

		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			msg := fmt.Sprintf("menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return 
		}
		c.JSON(http.StatusOK,result)

	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle :=context.WithTimeout(context.Background(),100*time.Second)
		defer cancle()
		var menu models.Menu
		menuId:=c.Param("menu_id")

		err :=c.BindJSON(&menu)
		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}
		validationErr :=validate.Struct(menu)
		if validationErr!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
				"error":validationErr.Error(),
			})
			return
		}
		

		filter :=bson.M{"menu_id":menuId}

		var updateObj primitive.D

		if menu.Start_date!=nil && menu.End_date !=nil{
			if !inTimeSpan(*menu.Start_date,*menu.End_date,time.Now()){
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError,gin.H{
					"error":msg,
				})
				return 
			}
			updateObj=append(updateObj, bson.E{Key: "start_date", Value: menu.Start_date})
			updateObj=append(updateObj, bson.E{Key: "end_date", Value: menu.End_date})

			if menu.Name !=""{
				updateObj =append(updateObj, bson.E{Key: "name",Value: menu.Name})
			}
			if menu.Category !=""{
				updateObj =append(updateObj, bson.E{Key: "category",Value: menu.Category})
			}

			menu.Updated_at=time.Now()
			updateObj=append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

			upsert :=true

			opt :=options.UpdateOptions{
				Upsert: &upsert,
			}

			result,err :=menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set",updateObj},
				},
				&opt,
			)

			if err!=nil{
				msg :="menu update failed "
				c.JSON(http.StatusInternalServerError,gin.H{
					"error":msg,
				})
				return 
			}
			c.JSON(http.StatusOK,result)
		}
	}
}
