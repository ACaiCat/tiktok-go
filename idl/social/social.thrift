namespace go social

include "model.thrift"
include "common.thrift"

// 关注操作类型
enum FollowActionType {
  // 关注
  FOLLOW = 1,
  // 取消关注
  UNFOLLOW = 2
}

// 关注请求
struct FollowReq {
  // 操作对象ID
  1: required string to_user_id (api.form = 'to_user_id');
  // 操作类型
  2: required FollowActionType action_type (api.form = 'action_type');
}

// 关注响应
struct FollowResp {
  // 响应状态
  1: required common.Base base;
}

// 社交用户数据
struct SocialUserData {
  // 社交用户列表
  1: required list<model.SocialUser> items;
  // 总数
  2: required i32 total;
}

// 关注列表请求
struct ListFollowingReq {
  // 用户ID
  1: required string user_id (api.query = 'user_id');
  // 页码
  2: required i32 page_num (api.query = 'page_num');
  // 单页尺寸
  3: required i32 page_size (api.query = 'page_size');
}

// 关注列表响应
struct ListFollowingResp {
  // 响应状态
  1: required common.Base base;
  // 响应数据
  2: optional SocialUserData data;
}

// 粉丝列表请求
struct ListFollowerReq {
  // 用户ID
  1: required string user_id (api.query = 'user_id');
  // 页码
  2: required i32 page_num (api.query = 'page_num');
  // 单页尺寸
  3: required i32 page_size (api.query = 'page_size');
}

// 粉丝列表响应
struct ListFollowerResp {
  // 响应状态
  1: required common.Base base;
  // 响应数据
  2: optional SocialUserData data;
}

// 好友列表请求
struct ListFriendReq {
  // 页码
  1: required i32 page_num (api.query = 'page_num');
  // 单页尺寸
  2: required i32 page_size (api.query = 'page_size');
}

// 好友列表响应
struct ListFriendResp {
  // 响应状态
  1: required common.Base base;
  // 响应数据
  2: optional SocialUserData data;
}

service SocialHandler {
  // 关注
  FollowResp Follow(1: FollowReq req) (api.post = "/relation/action")

  // 关注列表
  ListFollowingResp ListFollowing(1: ListFollowingReq req) (api.get = "/following/list")

  // 粉丝列表
  ListFollowerResp ListFollower(1: ListFollowerReq req) (api.get = "/follower/list")

  // 好友列表
  ListFriendResp ListFriend(1: ListFriendReq req) (api.get = "/friends/list")
}

