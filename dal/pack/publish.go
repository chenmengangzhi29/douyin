package pack

import (
	"github.com/chenmengangzhi29/douyin/dal/db"
	"github.com/chenmengangzhi29/douyin/kitex_gen/publish"
)

// VideoInfo pack video list info
func PublishInfo(currentId int64, videoData []*db.VideoRaw, userMap map[int64]*db.UserRaw, favoriteMap map[int64]*db.FavoriteRaw, relationMap map[int64]*db.RelationRaw) []*publish.Video {
	videoList := make([]*publish.Video, 0)
	for _, video := range videoData {
		videoUser, ok := userMap[video.UserId]
		if !ok {
			videoUser = &db.UserRaw{
				Name:          "未知用户",
				FollowCount:   0,
				FollowerCount: 0,
			}
			videoUser.ID = 0
		}

		var isFavorite bool = false
		var isFollow bool = false

		if currentId != -1 {
			_, ok := favoriteMap[int64(video.ID)]
			if ok {
				isFavorite = true
			}
			_, ok = relationMap[video.UserId]
			if ok {
				isFollow = true
			}
		}
		videoList = append(videoList, &publish.Video{
			Id: int64(video.ID),
			Author: &publish.User{
				Id:            int64(videoUser.ID),
				Name:          videoUser.Name,
				FollowCount:   videoUser.FollowCount,
				FollowerCount: videoUser.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		})
	}

	return videoList
}
