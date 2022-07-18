package main

import (
	"context"

	"github.com/chenmengangzhi29/douyin/cmd/relation/service"
	"github.com/chenmengangzhi29/douyin/dal/pack"
	"github.com/chenmengangzhi29/douyin/kitex_gen/relation"
	"github.com/chenmengangzhi29/douyin/pkg/errno"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	resp = new(relation.RelationActionResponse)

	if len(req.Token) == 0 || req.ToUserId == 0 || req.ActionType == 0 {
		resp.BaseResp = pack.BuilRelationBaseResp(errno.ParamErr)
		return resp, nil
	}

	err = service.NewRelationActionService(ctx).RelationAction(req)
	if err != nil {
		resp.BaseResp = pack.BuilRelationBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuilRelationBaseResp(errno.Success)
	return resp, nil
}

// FollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowList(ctx context.Context, req *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	resp = new(relation.FollowListResponse)

	if req.UserId == 0 {
		resp.BaseResp = pack.BuilRelationBaseResp(errno.ParamErr)
		return resp, nil
	}

	userList, err := service.NewFollowListService(ctx).FollowList(req)
	if err != nil {
		resp.BaseResp = pack.BuilRelationBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuilRelationBaseResp(errno.Success)
	resp.UserList = userList
	return resp, nil
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	resp = new(relation.FollowerListResponse)

	if req.UserId == 0 {
		resp.BaseResp = pack.BuilRelationBaseResp(errno.ParamErr)
		return resp, nil
	}

	userList, err := service.NewFollowerListService(ctx).FollowerList(req)
	if err != nil {
		resp.BaseResp = pack.BuilRelationBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuilRelationBaseResp(errno.Success)
	resp.UserList = userList
	return resp, nil
}
