package middleware

import (
	"aweme_kitex/pkg/logger"
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

// https://www.cloudwego.io/zh/docs/kitex/tutorials/framework-exten/middleware/

var _ endpoint.Middleware = ClientMiddleware

func ClientMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		ri := rpcinfo.GetRPCInfo(ctx)
		logger.Info("server address: %v, rpc timeout: %v, readwrite timeout: %v\n", ri.To().Address(), ri.Config().RPCTimeout(),
			ri.Config().ConnectTimeout())
		if err = next(ctx, req, resp); err != nil {
			return err
		}
		return nil
	}
}
