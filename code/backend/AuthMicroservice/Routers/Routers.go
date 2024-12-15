package Routers

import (
	"authmicroservice/Controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("login", Controllers.Login)
	r.POST("register", Controllers.Register)
	r.POST("logout", Controllers.Logout)
	r.DELETE("user", Controllers.Delete)
	r.GET("verify", Controllers.VerifyToken)

	return r
}