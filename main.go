package main

import (
	"douyin/model"
	"douyin/routers"
	"douyin/util"
	"os"
)

func main() {
	if err := Init(); err != nil {
		os.Exit(-1)
	}
	defer util.Logger.Sync()
	// 加载多个APP的路由配置
	routers.Include(routers.Video, routers.User, routers.Favorite, routers.Comment, routers.Follow)
	// 初始化路由
	r := routers.Init()
	if err := r.Run(); err != nil {
		return
	}
}

func Init() error {
	if err := util.InitLogger(); err != nil {
		return err
	}
	if err := model.Init(); err != nil {
		util.Logger.Error("model fail")
		return err
	}
	return nil
}
