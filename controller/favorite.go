package controller

import (
	"douyin/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	user_id := int(usersLoginInfo[token].Id)
	action := c.Query("action_type")
	favo := FavoriteRaw{}
	video := model.Videos{}
	model.DB.Table("favorite").Where("user_id=? AND video_id=?", user_id, video_id).First(&favo)
	if favo.Id == 0 {
		if action == "1" {
			model.DB.Table("video").Where("id=?", video_id).Find(&video)
			num := video.FavoriteCount + 1
			model.DB.Table("video").Model(&video).Where("id = ?", video_id).Update("favorite_count", num)
			last := FavoriteRaw{}
			model.DB.Table("favorite").Last(&last)
			model.DB.Table("favorite").Select("id", "user_id", "video_id").Create(&FavoriteRaw{Id: last.Id + 1, UserId: int64(user_id), VideoId: int64(video_id)})
		}
	} else {
		if action == "2" {
			model.DB.Table("video").Where("id=?", video_id).Find(&video)
			num := video.FavoriteCount - 1
			model.DB.Table("video").Where("id=?", video_id).Model(&video).Update("favorite_count", num)
			model.DB.Table("favorite").Delete(&favo)
		}

	}

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userid := c.Query("user_id")
	var videoIdList = make([]int, 10)
	fmt.Println(userid)
	model.DB.Select("video_id").Table("favorite").Where("user_id=?", userid).Find(&videoIdList)
	fmt.Println(videoIdList)
	var videosList []model.Videos

	var authorIdList = make([]int, 10)
	for i := 0; i < len(videoIdList); i++ {
		model.DB.Table("video").Select("id", "title", "play_url", "cover_url", "favorite_count", "comment_count").Where("id=?", videoIdList[i]).Find(&videosList)

	}
	model.DB.Table("video").Select("user_id").Where("id in ?", videoIdList).Find(&authorIdList)
	var videoList = make([]Video, len(videosList))
	for i := 0; i < len(videoList); i++ {
		videoList[i].Id = videosList[i].Id
		videoList[i].FavoriteCount = videosList[i].FavoriteCount
		videoList[i].PlayUrl = videosList[i].PlayUrl
		videoList[i].CoverUrl = videosList[i].CoverUrl
		videoList[i].CommentCount = videosList[i].CommentCount
		model.DB.Table("user").Select("id", "name", "follow_count", "follower_count").Find(&videoList[i].Author, authorIdList)
	}
	fmt.Println(videosList)
	fmt.Println(videoList)
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
