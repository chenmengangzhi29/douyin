package pack

import (
	"errors"
	"time"

	"github.com/chenmengangzhi29/douyin/kitex_gen/feed"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
)

//BuildBaseResp build baseResp from error
func BuildBaseResp(err error) *feed.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *feed.BaseResp {
	return &feed.BaseResp{StatusCode: err.ErrCode, StatusMessage: err.ErrMsg, ServiceTime: time.Now().Unix()}
}
