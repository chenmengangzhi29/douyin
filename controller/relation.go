package controller

import (
	"douyin/model"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}
type Relation struct {
	Id       int64 `gorm:"column:id"`
	UserID   int64 `gorm:"column:user_id"`
	ToUserId int64 `gorm:"column:to_user_id"`
	Status   int8  `gorm:"column:status"`
}

func (Relation) TableName() string {
	return "relation"
}

var db = model.DB

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	action_type, _ := strconv.Atoi(c.Query("action_type"))

	if action_type != 1 && action_type != 2 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "The value of action is invalid!"})
		return
	}

	if _, exist := usersLoginInfo[token]; exist {
		user_id := usersLoginInfo[token].Id
		var relation Relation
		db.Where("(user_id=? and to_user_id=?) or (user_id=? and to_user_id=?)", user_id, to_user_id, to_user_id, user_id).Find(&relation)
		fmt.Println(relation)
		fmt.Println(relation.Status)
		if !reflect.DeepEqual(relation, Relation{}) && relation.Status == 0 {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "The user has been followed!"})
			return
		}
		if reflect.DeepEqual(action_type, 1) {
			//follow
			if reflect.DeepEqual(relation, Relation{}) {
				fmt.Println("insert")
				new_relation := &Relation{
					UserID:   user_id,
					ToUserId: to_user_id,
					Status:   -1,
				}
				db.Create(&new_relation)
			} else {
				fmt.Println("update")
				var status int8 = 0
				if relation.UserID == user_id && relation.Status == 1 {
					status--
				} else if relation.UserID == to_user_id && relation.Status == -1 {
					status++
				} else {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Has been followed!"})
					return
				}
				db.Model(&relation).Update("status", relation.Status+status)
			}
		} else {
			// unfollow
			fmt.Println("unfollow")
			if (relation.UserID == user_id && relation.Status == -1) || (relation.ToUserId == user_id && relation.Status == 1) {
				db.Delete(&relation)
			} else {
				var status int8
				if relation.UserID == user_id {
					status++
				} else {
					status--
				}
				db.Model(&relation).Update("status", relation.Status+status)
			}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

type UserDao struct {
	Id            int64  `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	FollowCount   int64  `gorm:"column:follow_count"`
	FollowerCount int64  `gorm:"column:follower_count"`
}

func (UserDao) TableName() string {
	return "user"
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	user_id, _ := strconv.Atoi(c.Query("user_id"))

	if _, exist := usersLoginInfo[token]; exist {

		relation := make([]int, 0)
		FindFollowedUser(user_id, &relation)
		var userList = GetFollowedUserInfo(relation)

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: userList,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User does not exist!"},
		})
	}

}

//get the followed user info
func GetFollowedUserInfo(userlist []int) []User {
	user := make([]UserDao, 0)
	db.Table("user").Find(&user, map[string]interface{}{"id": userlist})

	response_user := make([]User, len(user))
	for i := 0; i < len(user); i++ {
		response_user[i].Id = user[i].Id
		response_user[i].Name = user[i].Name
		response_user[i].FollowCount = user[i].FollowCount
		response_user[i].FollowerCount = user[i].FollowerCount
		response_user[i].IsFollow = true
	}
	return response_user
}

// find out all followed user_id
func FindFollowedUser(user_id int, relation *[]int) {

	ch := make(chan []int, 2)
	go func() {
		var relation_a = make([]int, 0)
		db.Select("to_user_id").Table("relation").Where("user_id=? and status <= 0", user_id).Find(&relation_a)
		ch <- relation_a
	}()
	go func() {
		var relation_b = make([]int, 0)
		db.Select("user_id").Table("relation").Where("to_user_id=? and status >= 0", user_id).Find(&relation_b)
		ch <- relation_b
	}()

	*relation = append(<-ch, <-ch...)
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {

	token := c.Query("token")
	user_id, _ := strconv.Atoi(c.Query("user_id"))

	if _, exist := usersLoginInfo[token]; exist {

		fans := make([]int, 0)
		followed_fans := make([]int, 0)

		var wg sync.WaitGroup
		wg.Add(2)
		FindFansUser(user_id, &fans, &followed_fans, &wg)
		wg.Wait()

		//fmt.Println(fans, followed_fans)
		var userList = GetFansInfo(fans, followed_fans)

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: userList,
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

}

func GetFansInfo(fans []int, followed_fans []int) []User {

	user := make([]UserDao, 0)
	db.Table("user").Find(&user, map[string]interface{}{"id": fans})
	response_user := make([]User, len(user))
	for i := 0; i < len(user); i++ {
		response_user[i].Id = user[i].Id
		response_user[i].Name = user[i].Name
		response_user[i].FollowCount = user[i].FollowCount
		response_user[i].FollowerCount = user[i].FollowerCount
		response_user[i].IsFollow = false
	}

	db.Table("user").Find(&user, map[string]interface{}{"id": followed_fans})
	response_followed := make([]User, len(user))
	for i := 0; i < len(user); i++ {
		response_followed[i].Id = user[i].Id
		response_followed[i].Name = user[i].Name
		response_followed[i].FollowCount = user[i].FollowCount
		response_followed[i].FollowerCount = user[i].FollowerCount
		response_followed[i].IsFollow = true
	}
	response_user = append(response_user, response_followed...)
	return response_user
}

func FindFansUser(user_id int, fans *[]int, followed_fans *[]int, wg *sync.WaitGroup) {

	go func() {
		defer wg.Done()
		ch := make(chan []int, 2)
		var fan []int
		db.Select("to_user_id").Table("relation").Where("user_id=? and status = 0", user_id).Find(&fan)
		ch <- fan
		db.Select("user_id").Table("relation").Where("to_user_id=? and status = 0", user_id).Find(&fan)
		ch <- fan
		*fans = append(<-ch, <-ch...)

	}()

	go func() {
		defer wg.Done()
		ch := make(chan []int, 2)
		var followed = make([]int, 0)
		db.Select("user_id").Table("relation").Where("to_user_id=? and status = -1", user_id).Find(&followed)
		ch <- followed
		db.Select("to_user_id").Table("relation").Where("user_id=? and status = 1", user_id).Find(&followed)
		ch <- followed
		*followed_fans = append(<-ch, <-ch...)

	}()
}
