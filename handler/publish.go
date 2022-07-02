package handler

import (
	"douyin/model"
	"douyin/service"
	"mime/multipart"
	"strconv"
)

//---------------------handler--------------------------------
//处理传入参数，调用service层函数查询视频列表，封装响应

func PublishVideoData(token string, data multipart.File, title string) *model.Response {
	if length := len(token); length <= 0 || length > 64 {
		return &model.Response{
			StatusCode: -1, StatusMsg: "token length out of range(0,64]",
		}
	}

	if data == nil {
		return &model.Response{
			StatusCode: -1, StatusMsg: "file is empty",
		}
	}

	if length := len(title); length <= 0 || length > 128 {
		return &model.Response{
			StatusCode: -1, StatusMsg: "video title out of range(0,128]",
		}
	}

	err := service.PublishUserVideoData(token, data, title)
	if err != nil {
		return &model.Response{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	return &model.Response{
		StatusCode: 0, StatusMsg: "publish video success",
	}

}

func QueryVideoList(token string, userIdStr string) *model.VideoListResponse {
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return &model.VideoListResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	if len := len(token); len <= 0 || len > 64 {
		return &model.VideoListResponse{
			StatusCode: -1, StatusMsg: "token length out of range(0,64]",
		}
	}

	videoList, err := service.QueryUserVideoList(token, userId)
	if err != nil {
		return &model.VideoListResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		}
	}

	return &model.VideoListResponse{
		StatusCode: 0,
		StatusMsg:  "query video list success",
		VideoList:  videoList,
	}
}
