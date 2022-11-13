package routers

import (
	"aweme_kitex/controller"

	"github.com/gin-gonic/gin"
)

func Comment(apiRouter *gin.RouterGroup) {
	apiRouter.POST("/comments/action/", controller.CommentAction)
	apiRouter.GET("/comments/list/", controller.CommentList)
}
