package handler

import (
	"douyin/model"
	"douyin/service"
	"strconv"
)

//处理点赞操作的传入参数
func HandlerFavoriteAction(token string, videoIdStr string, actionTypeStr string) model.Response {
	if len := len(token); len <= 0 || len > 64 {
		return model.Response{StatusCode: -1, StatusMsg: "token length out of range(0,64]"}
	}

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		return model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		return model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	err = service.FavoriteActionData(token, videoId, actionType)
	if err != nil {
		return model.Response{StatusCode: -1, StatusMsg: err.Error()}
	}

	return model.Response{StatusCode: 0, StatusMsg: "favorite action success"}
}

//点赞列表的响应结构
type FavoriteListResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
}

//处理点赞列表的传入参数
func HandlerFavoriteList(userIdStr string, token string) FavoriteListResponse {
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return FavoriteListResponse{Response: model.Response{StatusCode: -1, StatusMsg: err.Error()}}
	}

	if len := len(token); len <= 0 || len > 64 {
		return FavoriteListResponse{Response: model.Response{StatusCode: -1, StatusMsg: "token length out of range(0,64]"}}
	}

	videoList, err := service.FavoriteListData(userId, token)
	if err != nil {
		return FavoriteListResponse{Response: model.Response{StatusCode: -1, StatusMsg: err.Error()}}
	}

	return FavoriteListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "favorite list success",
		},
		VideoList: videoList,
	}
}
