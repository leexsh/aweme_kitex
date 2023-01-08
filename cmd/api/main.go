package main

import (
	"aweme_kitex/cmd/api/handlers"
	"aweme_kitex/cmd/api/rpc"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/jwt"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/tracer"
	"go/constant"
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

	authMiddleware, err := jwt2.New(&jwt2.GinJWTMiddleware{
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
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVar handlers.UserLoginParam
			if err := c.ShouldBind(&loginVar); err != nil {
				return "", jwt2.ErrMissingLoginValues
			}
			if len(loginVar.UserName) == 0 || len(loginVar.Password) == 0 {
				return "", jwt2.ErrMissingLoginValues
			}
			// return rpc.CheckUser(context.Background(), &userService.CheckUserRequest{Username: loginVar.Username, Password: loginVar.Password})
			return nil, nil
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		logger.Fatal("JWT Error:" + err.Error())
	}

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt2.ExtractClaims(c)
		logger.Infof("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	aweme := r.Group("/aweme")
	aweme.GET("/feed/", handlers.Feed)

	publish := aweme.Group("/publish")
	publish.Use(authMiddleware.MiddlewareFunc())
	publish.POST("/action/", handlers.Publish)
	publish.GET("/list/", handlers.PublishList)

	user := aweme.Group("/user")
	user.GET("/", handlers.UserInfo)
	user.POST("/register/", handlers.Register)
	user.POST("/login/", authMiddleware.LoginHandler)

	favorite := aweme.Group("/favorite")
	favorite.Use(authMiddleware.MiddlewareFunc())
	favorite.POST("/action/", handlers.FavoriteAction)
	favorite.GET("/list/", handlers.FavoriteList)

	comment := aweme.Group("/comment")
	comment.Use(authMiddleware.MiddlewareFunc())
	comment.POST("/action/", handlers.CommentAction)
	comment.GET("/list/", handlers.CommentList)

	relation := aweme.Group("/relation")
	relation.POST("/action/", handlers.RelationAction)
	relation.GET("/follow/list/", handlers.FollowList)
	relation.GET("/follower/list/", handlers.FollowerList)

	if err := http.ListenAndServe(constant.ApiAddress, r); err != nil {
		logger.Fatal(err)
	}
}
