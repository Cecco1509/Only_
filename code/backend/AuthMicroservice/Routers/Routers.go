package Routers

import (
	"authmicroservice/Controllers"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware())

	r.POST("login", Controllers.Login)
	r.OPTIONS("login")
	r.POST("register", Controllers.Register)
	r.OPTIONS("register")
	r.POST("logout", Controllers.Logout)
	r.OPTIONS("logout")
	r.DELETE("user", Controllers.Delete)
	r.GET("verify", Controllers.VerifyToken)
	r.GET("userinfo/:username", Controllers.GetUserInfo)
	r.GET("users", Controllers.SearchUsers)

	return r
}