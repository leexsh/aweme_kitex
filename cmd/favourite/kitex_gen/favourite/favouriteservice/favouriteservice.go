// Code generated by Kitex v0.4.4. DO NOT EDIT.

package favouriteservice

import (
	favourite "aweme_kitex/cmd/favourite/kitex_gen/favourite"
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return favouriteServiceServiceInfo
}

var favouriteServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "FavouriteService"
	handlerType := (*favourite.FavouriteService)(nil)
	methods := map[string]kitex.MethodInfo{
		"FavouriteAction":       kitex.NewMethodInfo(favouriteActionHandler, newFavouriteServiceFavouriteActionArgs, newFavouriteServiceFavouriteActionResult, false),
		"FavouriteList":         kitex.NewMethodInfo(favouriteListHandler, newFavouriteServiceFavouriteListArgs, newFavouriteServiceFavouriteListResult, false),
		"QueryVideoIsFavourite": kitex.NewMethodInfo(queryVideoIsFavouriteHandler, newFavouriteServiceQueryVideoIsFavouriteArgs, newFavouriteServiceQueryVideoIsFavouriteResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "favourite",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.4.4",
		Extra:           extra,
	}
	return svcInfo
}

func favouriteActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*favourite.FavouriteServiceFavouriteActionArgs)
	realResult := result.(*favourite.FavouriteServiceFavouriteActionResult)
	success, err := handler.(favourite.FavouriteService).FavouriteAction(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFavouriteServiceFavouriteActionArgs() interface{} {
	return favourite.NewFavouriteServiceFavouriteActionArgs()
}

func newFavouriteServiceFavouriteActionResult() interface{} {
	return favourite.NewFavouriteServiceFavouriteActionResult()
}

func favouriteListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*favourite.FavouriteServiceFavouriteListArgs)
	realResult := result.(*favourite.FavouriteServiceFavouriteListResult)
	success, err := handler.(favourite.FavouriteService).FavouriteList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFavouriteServiceFavouriteListArgs() interface{} {
	return favourite.NewFavouriteServiceFavouriteListArgs()
}

func newFavouriteServiceFavouriteListResult() interface{} {
	return favourite.NewFavouriteServiceFavouriteListResult()
}

func queryVideoIsFavouriteHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*favourite.FavouriteServiceQueryVideoIsFavouriteArgs)
	realResult := result.(*favourite.FavouriteServiceQueryVideoIsFavouriteResult)
	success, err := handler.(favourite.FavouriteService).QueryVideoIsFavourite(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFavouriteServiceQueryVideoIsFavouriteArgs() interface{} {
	return favourite.NewFavouriteServiceQueryVideoIsFavouriteArgs()
}

func newFavouriteServiceQueryVideoIsFavouriteResult() interface{} {
	return favourite.NewFavouriteServiceQueryVideoIsFavouriteResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) FavouriteAction(ctx context.Context, req *favourite.FavouriteActionRequest) (r *favourite.FavouriteActionResponse, err error) {
	var _args favourite.FavouriteServiceFavouriteActionArgs
	_args.Req = req
	var _result favourite.FavouriteServiceFavouriteActionResult
	if err = p.c.Call(ctx, "FavouriteAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) FavouriteList(ctx context.Context, req *favourite.FavouriteListRequest) (r *favourite.FavouriteListResponse, err error) {
	var _args favourite.FavouriteServiceFavouriteListArgs
	_args.Req = req
	var _result favourite.FavouriteServiceFavouriteListResult
	if err = p.c.Call(ctx, "FavouriteList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) QueryVideoIsFavourite(ctx context.Context, req *favourite.QueryVideoIsFavouriteRequest) (r *favourite.QueryVideoIsFavouriteResponse, err error) {
	var _args favourite.FavouriteServiceQueryVideoIsFavouriteArgs
	_args.Req = req
	var _result favourite.FavouriteServiceQueryVideoIsFavouriteResult
	if err = p.c.Call(ctx, "QueryVideoIsFavourite", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
