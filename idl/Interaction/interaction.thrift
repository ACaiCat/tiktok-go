namespace go interaction

include "model.thrift"
include "common.thrift"

// 点赞请求
struct LikeReq {
  // 视频ID
  1: optional string video_id (api.form = 'video_id');
  // 评论ID
  2: optional string comment_id (api.form = 'comment_id');
  // 操作类型
  3: required i32 action_type (api.form = 'action_type');
}

// 点赞响应
struct LikeResp {
  // 响应状态
  1: required common.Base base;
}

// 视频数据
struct VideoData {
  // 视频列表
  1: required list<model.Video> items;
}

// 点赞列表请求
struct ListLikeReq {
  // 用户ID
  1: required string user_id (api.query = 'user_id');
  // 页码
  2: required string page_num (api.query = 'page_num');
  // 单页尺寸
  3: required string page_size (api.query = 'page_size');
}

// 点赞列表响应
struct ListLikeResp {
  // 响应状态
  1: required common.Base base;
  // 响应数据
  2: optional VideoData data;
}

// 评论请求
struct CommentReq {
  // 视频ID
  1: optional string video_id (api.form = 'video_id');
  // 评论ID
  2: optional string comment_id (api.form = 'comment_id');
  // 评论内容
  3: required string content (api.form = 'content');
}

// 评论响应
struct CommentResp {
  // 响应状态
  1: required common.Base base;
}

// 评论数据
struct CommentData {
  // 评论列表
  1: required list<model.Comment> items;
}

// 评论列表请求
struct ListCommentReq {
  // 视频ID
  1: optional string video_id (api.query = 'video_id');
  // 评论ID
  2: optional string comment_id (api.query = 'comment_id');
  // 页码
  3: required string page_num (api.query = 'page_num');
  // 单页尺寸
  4: required string page_size (api.query = 'page_size');
}

// 评论列表响应
struct ListCommentResp {
  // 响应状态
  1: required common.Base base;
  // 响应数据
  2: optional CommentData data;
}

// 删除评论请求
struct DeleteCommentReq {
  // 评论ID
  1: required string comment_id (api.form = 'comment_id');
}

// 删除评论响应
struct DeleteCommentResp {
  // 响应状态
  1: required common.Base base;
}

service InteractionHandler {
  // 点赞操作
  LikeResp Like(1: LikeReq req) (api.post = "/like/action")

  // 点赞列表
  ListLikeResp ListLike(1: ListLikeReq req) (api.get = "/like/list")

  // 评论
  CommentResp Comment(1: CommentReq req) (api.post = "/comment/publish")

  // 评论列表
  ListCommentResp ListComment(1: ListCommentReq req) (api.get = "/comment/list")

  // 删除评论
  DeleteCommentResp DeleteComment(1: DeleteCommentReq req) (api.delete = "/comment/delete")
}

