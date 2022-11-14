package routers

import (
	"aweme_kitex/controller"

	"github.com/gin-gonic/gin"
)

func Follow(apiRouter *gin.RouterGroup) {
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
