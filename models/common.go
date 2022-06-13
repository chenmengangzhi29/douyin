package models

//共享结构
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	VideoId       int64  `json:"video_id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type Comment struct {
	CommentId  int64  `json:"comment_id,omitempty"`
	UserId     int64  `json:"user_id,omitempty"`
	VideoId    int64  `json:"video_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	UserId        int64  `json:"user_id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type Relation struct {
	Id       int64 `json:"id,omitempty"`
	UserID   int64 `json:"user_id,omitempty"`
	ToUserId int64 `json:"to_user_id,omitempty"`
	Status   byte  `json:"status,omitempty"`
}

type Users struct{
	Id				int64`json:"id"`
	Name 			string`json:"name"`
	Password		 string`json:"password"`
	FanNum   		int64`json:"fan_num"`
	FollowNum		int64`json:"follow_num"`
}
type Videos struct {
	Id            int64  `json:"id,omitempty"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	Create_time   time.Time `json:"create___time"`
	Title string `json:"title"`
}

