package handler

import (
	"douyin/model"
	"douyin/service"
	"strconv"
)

//----------------------------handler-------------------------------------
//该层功能包括处理传入参数，向service层获取视频信息，封装成响应信息
type FeedResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

func QueryVideoFeed(token string, latestTimeStr string) FeedResponse {
	//处理传入参数
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		return FeedResponse{
			Response: model.Response{StatusCode: -1, StatusMsg: err.Error()},
		}
	}

	if len := len(token); len <= 0 || len > 64 {
		return FeedResponse{
			Response: model.Response{StatusCode: -1, StatusMsg: "token length out of range(0,64]"},
		}
	}

	//获取视频
	videoList, nextTime, err := service.QueryVideoData(latestTime, token)
	if err != nil {
		return FeedResponse{
			Response: model.Response{StatusCode: -1, StatusMsg: err.Error()},
		}
	}

	// fmt.Println(videoData)

	return FeedResponse{
		Response:  model.Response{StatusCode: 0, StatusMsg: "query video feed success"},
		VideoList: videoList,
		NextTime:  nextTime,
	}
}
