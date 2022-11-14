package routers

import "github.com/gin-gonic/gin"

type Option func(group *gin.RouterGroup)

var options = []Option{}

func Include(opts ...Option) {
	options = append(options, opts...)
}
func Init() *gin.Engine {
	r := gin.New()
	r.Static("./static", "./public")
	apiRouter := r.Group("/aweme")
	for _, opt := range options {
		opt(apiRouter)
	}
	return r
}
