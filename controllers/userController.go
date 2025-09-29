package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parthav/restraunt-management/database"
	"github.com/parthav/restraunt-management/helpers"
	"github.com/parthav/restraunt-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	
	"go.mongodb.org/mongo-driver/mongo"
)
var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		recordPerPage, err :=strconv.Atoi(c.Query("recordPerPage"))

		if err !=nil ||recordPerPage<1 {
			recordPerPage=10
		}
		page,err1:=strconv.Atoi(c.Query("page"))
		if err1 !=nil ||page<1 {
			page=1
		}
		startIndex:=(page-1)*recordPerPage

		matchStage :=bson.D{{"$match", bson.D{}}}
		 projectStage := bson.D{
            {"$project", bson.D{
                {"_id", 0},
                {"total_count", 1},
                {"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
            }},
        }
		result, err :=userCollection.Aggregate(ctx,
												mongo.Pipeline{
													matchStage,projectStage})

		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error while listing users",
			})
			return 
		}		
		var allUsers []bson.M
		err =result.All(ctx,&allUsers)		
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error while listing users",
			})
			return 
		}			
		c.JSON(http.StatusOK,allUsers)					
			}
			
			}
		




				

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		 ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var userId =c.Param("user_id")
		var user models.User

		err :=userCollection.FindOne(ctx,bson.M{"user_id":userId}).Decode(&user)
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error retriving user from db",
			})
			return 
		}
		c.JSON(http.StatusOK,user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var user models.User
		err :=c.BindJSON(&user)
		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}
		validErr :=validate.Struct(user)
		if validErr !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":validErr.Error(),
			})
			return 
		}
		count, err :=userCollection.CountDocuments(ctx,bson.M{"email":user.Email})
		if err !=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
			})
			return 
		}
		
		password :=HashedPassword(*user.Password)
		user.Password=&password

		count, err =userCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		if err !=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
			})
			return 
		}
		if count >0{
			c.JSON(409,gin.H{
				"error":"this email/phone already exists",
			})
			return 
		}
		user.Created_at=time.Now()
		user.Updated_at=time.Now()
		user.ID=primitive.NewObjectID()
		user.User_id=user.ID.Hex()

		token, refreshToken, _ :=helpers.GenerateAllTokens(*user.Email,*user.First_name,*user.Last_name,user.User_id)
		user.Token=&token
		user.Refresh_token=&refreshToken

		result, err :=userCollection.InsertOne(ctx,&user)

		if err !=nil{
		c.JSON(500,gin.H{
						"error":"error inserting user to db",
					})
					return 
		}
		c.JSON(http.StatusOK,result)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx , cancel :=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var user models.User
		 var foundUser models.User

		 err :=c.BindJSON(&user)
		 if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		 }
		 err =userCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&foundUser)
		 if err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{
						"error":"user not found",
					})
					return 
		 }
		 passwordIsValid, msg :=VerifyPassword(*user.Password,*foundUser.Password)

		 if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return 
		 }

		 token,refreshToken, _ :=helpers.GenerateAllTokens(*foundUser.Email,*foundUser.First_name,*foundUser.Last_name,foundUser.User_id)
		 
		 helpers.UpdateAllTokens(token,refreshToken,foundUser.User_id)
		 c.JSON(http.StatusOK,foundUser)


	}
}

func HashedPassword(password string) string {
	bytes,err :=bcrypt.GenerateFromPassword([]byte(password),14)
	if err !=nil{
		log.Panic(err)
	}
	return string(bytes)

}

func VerifyPassword(userpassword string, providedPassword string) (bool, string) {
	err :=bcrypt.CompareHashAndPassword([]byte(providedPassword),[]byte(userpassword))
	check :=true
	msg :=""
	if err !=nil{
		log.Panic(err)
		return false,"login or password is incorrect"
	}
	return check,msg

}
