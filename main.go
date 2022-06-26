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
	// 加载多个APP的路由配置
	routers.Include(routers.Video, routers.User, routers.Favorite, routers.Comment, routers.Follow)
	// 初始化路由
	r := routers.Init()
	// r.Use(gin.Logger())
	if err := r.Run(); err != nil {
		return
	}
}

func Init() error {
	if err := model.Init(); err != nil {
		return err
	}
	if err := util.InitLogger(); err != nil {
		return err
	}
	return nil
}
