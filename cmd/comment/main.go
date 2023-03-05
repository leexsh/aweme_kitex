package main

import (
	"aweme_kitex/cfg"
	comment "aweme_kitex/cmd/comment/kitex_gen/comment/commentservice"
	"aweme_kitex/pkg/bound"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/middleware"
	"aweme_kitex/pkg/tracer"
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"github.com/sirupsen/logrus"
)

func Init() {
	cfg.Init()
	logger.DoInit("", "comment_log", logrus.DebugLevel)
	tracer.InitJaeger(constants.FeedServiceName)
}

func main() {
	// 注册中心
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress}) // r should not be reused.
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", constants.CommentAddress)
	if err != nil {
		panic(err)
	}
	Init()
	svr := comment.NewServer(new(CommentServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.CommentServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr),                                       // address
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		server.WithSuite(trace.NewDefaultServerSuite()),                    // tracer
		server.WithBoundHandler(bound.NewCpuLimitHandler()),                // BoundHandler
		server.WithRegistry(r),                                             // registry
	)

	// svr := feed.NewServer(new(FeedServiceImpl))

	err = svr.Run()

	if err != nil {
		logger.Info(err.Error())
	}
}
