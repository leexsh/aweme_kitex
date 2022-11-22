package routers

import (
	"aweme_kitex/controller"

	"github.com/gin-gonic/gin"
)

func Feed(apiRouter *gin.RouterGroup) {
	apiRouter.POST("/feed", controller.Feed)
}
