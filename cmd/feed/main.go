package main

import (
	feed "aweme_kitex/cmd/feed/kitex_gen/feed/feedservice"
	videoKafka "aweme_kitex/cmd/feed/service_feed/kafka"
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
	videoKafka.InitKafka()
	logger.DoInit("", "feed_log", logrus.DebugLevel)
	tracer.InitJaeger(constants.FeedServiceName)
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress}) // r should not be reused.
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", constants.FeedAddress)
	if err != nil {
		panic(err)
	}
	Init()
	svr := feed.NewServer(new(FeedServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.FeedServiceName}), // server name
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
