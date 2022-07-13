// Code generated by Kitex v0.3.2. DO NOT EDIT.
package relationservice

import (
	"github.com/chenmengangzhi29/douyin/kitex_gen/relation"
	"github.com/cloudwego/kitex/server"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler relation.RelationService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}
