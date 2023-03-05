package main

import (
	"aweme_kitex/cfg"
	user "aweme_kitex/cmd/user/kitex_gen/user/userservice"
	userRPC "aweme_kitex/cmd/user/rpc"
	"aweme_kitex/cmd/user/service_user/db"
	userKafka "aweme_kitex/cmd/user/service_user/kafka"
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
	// redis init
	db.InitRedis()
	// kafka init
	userKafka.InitKafka()
	// rpc init
	userRPC.Init()
	cfg.Init()
	logger.DoInit("", "user_log", logrus.DebugLevel)
	tracer.InitJaeger(constants.UserServiceName)
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress}) // r should not be reused.
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", constants.UserAddress)
	if err != nil {
		panic(err)
	}
	Init()
	svr := user.NewServer(new(UserServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.UserServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr),                                       // address
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		server.WithSuite(trace.NewDefaultServerSuite()),                    // tracer
		server.WithBoundHandler(bound.NewCpuLimitHandler()),                // BoundHandler
		server.WithRegistry(r),                                             // registry
	)

	err = svr.Run()

	if err != nil {
		logger.Info(err.Error())
	}
}