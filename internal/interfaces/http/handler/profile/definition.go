package profile

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "github.com/gin-gonic/gin"

type ProfileHandler interface {
	GetMyProfile(c *gin.Context)
	UpdateMyProfile(c *gin.Context)
	UploadAvatar(c *gin.Context)
	GetPublicProfile(c *gin.Context)
}
