package middleware

import (
	"context"

	"github.com/chenmengangzhi29/douyin/pkg/logger"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

var _ endpoint.Middleware = CommonMiddleware

//CommonMiddleware common middleware print some rpc info, real request and real response
func CommonMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		ri := rpcinfo.GetRPCInfo(ctx)
		//get real request
		logger.Infof("real request: %+v\n", req)
		//get remote service information
		logger.Infof("remote service name: %s, remote method: %s\n", ri.To().ServiceName(), ri.To().Method())
		if err = next(ctx, req, resp); err != nil {
			return err
		}
		//get real response
		logger.Infof("real response: %+v\n", resp)
		return nil
	}
}
