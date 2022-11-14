package routers

import (
	"aweme_kitex/controller"

	"github.com/gin-gonic/gin"
)

func Favourite(apiRouter *gin.RouterGroup) {
	apiRouter.POST("/favourite/action/", controller.FavouriteAction)
	apiRouter.GET("/favourite/list/", controller.FavouriteList)
}
