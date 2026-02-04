package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterProfileRoutes(v1 *gin.RouterGroup, p Params, sessionAuth gin.HandlerFunc) {
	profile := v1.Group("/profile", sessionAuth)
	{
		profile.GET("", p.ProfileHandler.GetMyProfile)
		profile.PUT("", p.ProfileHandler.UpdateMyProfile)
		profile.POST("/avatar", p.ProfileHandler.UploadAvatar)
	}
}
