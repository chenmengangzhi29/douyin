package model

//共享结构

//-----------响应结构---------------
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FeedResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,omitempty"`
	NextTime   int64       `json:"next_time,omitempty"`
}

type VideoListResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,omitempty"`
}

type UserResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id,omitempty"`
	Token      string `json:"token,omitempty"`
}

type UserInfoResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	User       interface{} `json:"user"`
}

type FavoriteListResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	VideoList  interface{} `json:"video_list,omitempty"`
}

type CommentActionResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	Comment    interface{} `json:"comment,omitempty"`
}

type CommentListResponse struct {
	StatusCode  int32       `json:"status_code"`
	StatusMsg   string      `json:"status_msg"`
	CommentList interface{} `json:"comment_list,omitempty"`
}

type RelationListResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	UserList   interface{} `json:"user_list,omitempty"`
}

type Video struct {
	Id            int64  `json:"id"`
	Author        *User  `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type Comment struct {
	Id         int64  `json:"id"`
	User       *User  `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

//-----------------直接从数据库取数据的结构----------------------------
type VideoRaw struct {
	Id            int64  `gorm:"column:id"`
	UserId        int64  `gorm:"column:user_id"`
	Title         string `gorm:"column:title"`
	PlayUrl       string `gorm:"column:play_url"`
	CoverUrl      string `gorm:"column:cover_url"`
	FavoriteCount int64  `gorm:"column:favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count"`
	CreateTime    int64  `gorm:"column:create_time"`
}

type UserRaw struct {
	Id            int64  `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	Password      string `gorm:"column:password"`
	FollowCount   int64  `gorm:"column:follow_count"`
	FollowerCount int64  `gorm:"column:follower_count"`
	Token         string `gorm:"column:token"`
}

type FavoriteRaw struct {
	Id      int64 `gorm:"column:id"`
	UserId  int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
}

type RelationRaw struct {
	Id       int64 `gorm:"column:id"`
	UserId   int64 `gorm:"column:user_id"`
	ToUserId int64 `gorm:"column:to_user_id"`
}

type CommentRaw struct {
	Id         int64  `gorm:"column:id"`
	UserId     int64  `gorm:"column:user_id"`
	VideoId    int64  `gorm:"column:video_id"`
	Contents   string `grom:"column:contents"`
	CreateDate string `grom:"column:create_date"`
}
