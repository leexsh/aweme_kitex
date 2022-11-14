package routers

import (
	"aweme_kitex/controller"

	"github.com/gin-gonic/gin"
)

func Video(apiRouter *gin.RouterGroup) {
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/publish/action", controller.Publish)
	apiRouter.GET("/publish/list", controller.PublishList)

}
