package main

import (
	"douyin/model"
	"douyin/routers"
	"douyin/util/logger"
	"os"
)

func main() {
	if err := Init(); err != nil {
		os.Exit(-1)
	}
	defer logger.Sync()
	// 加载多个APP的路由配置
	routers.Include(routers.Video, routers.User, routers.Favorite, routers.Comment, routers.Follow)
	// 初始化路由
	r := routers.Init()
	if err := r.Run(); err != nil {
		return
	}
}

func Init() error {
	if err := logger.Init(); err != nil {
		return err
	}
	if err := model.ConfigInit(); err != nil {
		return err
	}
	if err := model.MysqlInit(); err != nil {
		return err
	}
	if err := model.OssInit(); err != nil {
		return err
	}
	return nil
}
