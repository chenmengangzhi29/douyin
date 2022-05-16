package main

import (
	"douyin/routers"
	"fmt"
)

func main() {
	// 加载多个APP的路由配置
	routers.Include(routers.Video, routers.User, routers.Favorite, routers.Comment, routers.Follow)
	// 初始化路由
	r := routers.Init()
	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
