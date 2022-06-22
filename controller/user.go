package controller

import (
	"douyin/model"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	//"sync/atomic"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}
var u = model.Users{}
var last = model.Users{}
var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password
	model.DB.Where("name = ?", username).First(&u)
	if _, exist := usersLoginInfo[token]; exist || u.Name == username {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Id:   userIdSequence,
			Name: username,
		}
		usersLoginInfo[token] = newUser
		model.DB.Last(&last)
		newUsers := model.Users{Id: last.Id + 1, Name: username, Password: password, FanNum: 0, FollowNum: 0}
		model.DB.Select("id", "name", "password", "fan_num", "follow_num").Create(&newUsers)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}
	//if u.Name==username{
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	//	})
	//} else {
	//	models.DB.Last(&last)
	//	newUsers :=models.Users{Id: last.Id+1,Name: username,Password: password,FanNum: 0,FollowNum: 0}
	//	models.DB.Select("id","name","password","fan_num","follow_num").Create(&newUsers)
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   newUsers.Id,
	//		Token:    username + password,
	//	})
	//}

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password
	model.DB.Where("name = ?", username).Find(&u)
	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else if username == u.Name {
		token = u.Name + u.Password
		usersLoginInfo[token] = User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowNum,
			FollowerCount: u.FanNum,
			IsFollow:      false,
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist "},
		})
	}
	//if token==u.Name+u.Password {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   u.Id,
	//		Token:    token,
	//	})
	//}else if username==u.Name{
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   u.Id,
	//		Token:    token,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist1 "},
	//	})
	//}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else if token == u.Name+u.Password {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: User{
				Id:            u.Id,
				Name:          u.Name,
				FollowCount:   u.FollowNum,
				FollowerCount: u.FanNum,
				IsFollow:      true,
			},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist "},
		})
	}
	//if token==u.Name+u.Password {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 0},
	//		User:     User{
	//			Id: u.Id,
	//			Name: u.Name,
	//			FollowCount: u.FollowNum,
	//			FollowerCount: u.FanNum,
	//			IsFollow: true,
	//		},
	//	})
	//}else {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist "},
	//	})
	//}
}
