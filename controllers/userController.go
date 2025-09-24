package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func Login() gin.HandlerFunc{
	return func(ctx *gin.Context){

	}
}

func HashedPassword(password string) string{

}

func VerifyPassword(userpassword string, providedPassword string)(bool,string){
	
}