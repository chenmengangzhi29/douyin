package pack

import (
	"errors"
	"time"

	"github.com/chenmengangzhi29/douyin/kitex_gen/feed"
	"github.com/chenmengangzhi29/douyin/kitex_gen/publish"
	"github.com/chenmengangzhi29/douyin/kitex_gen/user"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
)

//BuildFeedBaseResp build feed baseResp from error
func BuildFeedBaseResp(err error) *feed.BaseResp {
	if err == nil {
		return feedbaseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return feedbaseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return feedbaseResp(s)
}

func feedbaseResp(err errno.ErrNo) *feed.BaseResp {
	return &feed.BaseResp{StatusCode: err.ErrCode, StatusMessage: err.ErrMsg, ServiceTime: time.Now().Unix()}
}

//BuildPublishBaseResp build publish baseResp from error
func BuildPublishBaseResp(err error) *publish.BaseResp {
	if err == nil {
		return publishbaseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return publishbaseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return publishbaseResp(s)
}

func publishbaseResp(err errno.ErrNo) *publish.BaseResp {
	return &publish.BaseResp{StatusCode: err.ErrCode, StatusMessage: err.ErrMsg, ServiceTime: time.Now().Unix()}
}

//BuildUserBaseResp build user baseResp from error
func BuildUserBaseResp(err error) *user.BaseResp {
	if err == nil {
		return userbaseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return userbaseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return userbaseResp(s)
}

func userbaseResp(err errno.ErrNo) *user.BaseResp {
	return &user.BaseResp{StatusCode: err.ErrCode, StatusMessage: err.ErrMsg, ServiceTime: time.Now().Unix()}
}
