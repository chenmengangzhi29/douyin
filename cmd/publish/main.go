package main

import (
	"net"

	"github.com/chenmengangzhi29/douyin/dal"
	publish "github.com/chenmengangzhi29/douyin/kitex_gen/publish/publishservice"
	"github.com/chenmengangzhi29/douyin/pkg/bound"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	"github.com/chenmengangzhi29/douyin/pkg/middleware"
	"github.com/chenmengangzhi29/douyin/pkg/oss"
	tracer2 "github.com/chenmengangzhi29/douyin/pkg/tracer"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

func Init() {
	tracer2.InitJaeger(constants.PublishServiceName)
	dal.Init()
	oss.Init()
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress}) // r should not be reused
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", constants.PublishAddress)
	if err != nil {
		panic(err)
	}
	Init()
	svr := publish.NewServer(new(PublishServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.PublishServiceName}), //server name                                               // middleWare
		server.WithMiddleware(middleware.ServerMiddleware),                                                // middleWare
		server.WithServiceAddr(addr),                                                                      //address
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),                                //limit
		server.WithMuxTransport(),                                                                         //Multiplex
		server.WithSuite(trace.NewDefaultServerSuite()),                                                   //tracer
		server.WithBoundHandler(bound.NewCpuLimitHandler()),                                               //BoundHandler
		server.WithRegistry(r),                                                                            // registry
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
