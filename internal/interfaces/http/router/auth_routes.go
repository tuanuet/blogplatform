package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(v1 *gin.RouterGroup, p Params, rateLimit, sessionAuth gin.HandlerFunc) {
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", rateLimit, p.AuthHandler.Register)
		authGroup.POST("/login", rateLimit, p.AuthHandler.Login)
		authGroup.POST("/logout", sessionAuth, p.AuthHandler.Logout)
		authGroup.GET("/:provider", p.AuthHandler.SocialLogin)
		authGroup.GET("/:provider/callback", p.AuthHandler.SocialCallback)
	}
}
