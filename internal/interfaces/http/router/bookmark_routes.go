package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterBookmarkRoutes(v1 *gin.RouterGroup, p Params, sessionAuth gin.HandlerFunc) {
	v1.GET("/bookmarks", sessionAuth, p.BookmarkHandler.List)
}
