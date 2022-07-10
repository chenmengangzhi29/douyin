package rpc

import (
	"github.com/chenmengangzhi29/douyin/kitex_gen/feed/feedservice"
	"github.com/chenmengangzhi29/douyin/pkg/constants"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var feedClient feedservice.Client

func initFeedRpc() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

}
