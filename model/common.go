package model

//共享结构

//响应结构
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type Video struct {
	Id            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type Comment struct {
	Id         int64  `json:"id"`
	User       User   `json:"user"`
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

//直接从数据库取数据的结构
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

// //其他
// type Users struct {
// 	Id        int64  `json:"id"`
// 	Name      string `json:"name"`
// 	Password  string `json:"password"`
// 	FanNum    int64  `json:"fan_num"`
// 	FollowNum int64  `json:"follow_num"`
// }
// type Videos struct {
// 	Id            int64     `json:"id,omitempty"`
// 	PlayUrl       string    `json:"play_url" json:"play_url,omitempty"`
// 	CoverUrl      string    `json:"cover_url,omitempty"`
// 	FavoriteCount int64     `json:"favorite_count,omitempty"`
// 	CommentCount  int64     `json:"comment_count,omitempty"`
// 	Create_time   time.Time `json:"create___time"`
// 	Title         string    `json:"title"`
// }
