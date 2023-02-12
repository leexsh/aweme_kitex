package main

import (
	"aweme_kitex/cmd/api/handlers"
	"aweme_kitex/cmd/api/rpc"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/errno"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/tracer"
	"net/http"
	"time"

	jwt2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func Init() {
	tracer.InitJaeger(constants.ApiServiceName)
	rpc.InitRPC()
}

func main() {
	Init()
	r := gin.New()

	authMiddleware, _ := jwt2.New(&jwt2.GinJWTMiddleware{
		Key:     []byte(jwt.JwtKey),
		Timeout: time.Hour * 24 * 365 * 10,
		PayloadFunc: func(data interface{}) jwt2.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt2.MapClaims{
					constants.IdentiryKey: v,
				}
			}
			return jwt2.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			uc, err := jwt.AnalyzeToken(message)
			if err != nil {
				panic(err)
			}
			handlers.SendUserResponse(c, errno.Success, uc.Id, message)
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	aweme := r.Group("/aweme")
	aweme.GET("/feed/", handlers.Feed)

	publish := aweme.Group("/publish")
	publish.POST("/action/", handlers.Publish)
	publish.GET("/list/", handlers.PublishList)

	user := aweme.Group("/user")
	user.GET("/", handlers.UserInfo)
	user.POST("/register/", handlers.Register)
	user.POST("/login/", authMiddleware.LoginHandler)

	favorite := aweme.Group("/favorite")
	favorite.POST("/action/", handlers.FavoriteAction)
	favorite.GET("/list/", handlers.FavoriteList)

	comment := aweme.Group("/comment")
	comment.POST("/action/", handlers.CommentAction)
	comment.GET("/list/", handlers.CommentList)

	relation := aweme.Group("/relation")
	relation.POST("/action/", handlers.RelationAction)
	relation.GET("/follow/list/", handlers.FollowList)
	relation.GET("/follower/list/", handlers.FollowerList)

	if err := http.ListenAndServe(constants.ApiAddress, r); err != nil {
		logger.Fatal(err)
	}
}
