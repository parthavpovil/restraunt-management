package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthav/restraunt-management/helpers"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context) {
		clientToken :=c.Request.Header.Get("token")
		if clientToken==""{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"no auth token provided",
			})
			c.Abort()
			return 
		}

		claims, err :=helpers.ValidateToken(clientToken)
		if err!=""{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":err,
			})
			c.Abort()
			return 
		}
		c.Set("email",claims.Email)
		c.Set("first_name",claims.First_name)
		c.Set("last_name",claims.Last_name)
		c.Set("uid",claims.Uid)
		c.Next()

	}
}