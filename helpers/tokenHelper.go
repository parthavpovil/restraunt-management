package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/parthav/restraunt-management/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstname string, lastname string, uid string) (signedTOken string, signedRefreshTokenrefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstname,
		Last_name:  lastname,
		Uid:        uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(12))),
		},
	}
	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}
	token,err :=jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err :=jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims).SignedString([]byte(SECRET_KEY))
	if err !=nil{
		log.Fatal(err)
		return
	}
	return token,refreshToken,err
}

func UpdateAllTokens(signedToken string,signedRefreshToken string, userId string) {
	var ctx, cancel =context.WithTimeout(context.Background(),100*time.Second)
	var updateObj primitive.D
	updateObj=append(updateObj, bson.E{"token",signedToken})
	updateObj=append(updateObj,bson.E{"refresh_token",signedRefreshToken })
	updated_at :=time.Now()

	updateObj=append(updateObj, bson.E{"updated_at",updated_at})

	upsert :=true

	filter:=bson.M{"user_id":userId}
	opt :=options.UpdateOptions{
		Upsert: &upsert,
	}
	 _,err :=userCollection.UpdateOne(ctx,
		filter,
		bson.D{
			{"$set",updateObj},
		},
		&opt,
	)
	defer cancel()
	if err !=nil{
		log.Panic(err)
		return
	}
	return


}

func ValidateToken(signedToken string)(claims *SignedDetails, msg string) {
	jwt.ParseWithClaims()
}
