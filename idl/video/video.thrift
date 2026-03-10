namespace go video

include "model.thrift"
include "common.thrift"

// 视频流请求
struct FeedReq {
  // 时间戳游标
  1: optional string latestTime (api.query = 'latest_time');
}

// 视频流响应
struct FeedResp {
  // 响应状态
  1: required common.Base base;
  // 视频列表
  2: optional list<model.Video> items;
}

// 投稿视频请求
struct PublishReq {
  // 视频文件
  1: required binary data (api.form = 'data');
  // 视频标题
  2: required string title (api.form = 'title');
  // 视频描述
  3: required string description (api.form = 'description');
}

// 投稿视频响应
struct PublishResp {
  // 响应状态
  1: required common.Base base;
}

// 发布列表请求
struct ListReq {
  // 用户ID
  1: required string userID (api.query = 'user_id');
  // 页码
  2: required string pageNum (api.query = 'page_num');
  // 单页尺寸
  3: required string pageSize (api.query = 'page_size');
}

// 发布列表响应
struct ListResp {
  // 响应状态
  1: required common.Base base;
  // 视频列表
  2: optional list<model.Video> items;
}

// 热门排行榜请求
struct PopularReq {
  // 页码
  1: required string pageNum (api.query = 'page_num');
  // 单页尺寸
  2: required string pageSize (api.query = 'page_size');
}

// 热门排行榜响应
struct PopularResp {
  // 响应状态
  1: required common.Base base;
  // 视频列表
  2: optional list<model.Video> items;
}

// 搜索视频请求
struct SearchReq {
  // 关键词
  1: required string keywords (api.form = 'keywords');
  // 页码
  2: required string pageNum (api.form = 'page_num');
  // 单页尺寸
  3: required string pageSize (api.form = 'page_size');
  // 起始时间
  4: optional string fromDate (api.form = 'from_date');
  // 结束时间
  5: optional string toDate (api.form = 'to_date');
  // 用户名关键词
  6: optional string username (api.form = 'username');
}

// 搜索视频响应
struct SearchResp {
  // 响应状态
  1: required common.Base base;
  // 视频列表
  2: optional list<model.Video> items;
}

service VideoHandler {
  // 视频流
  FeedResp Feed(1: FeedReq req) (api.get = "/video/feed")

  // 投稿
  PublishResp Publish(1: PublishReq req) (api.post = "/video/publish")

  // 发布列表
  ListResp List(1: ListReq req) (api.get = "/video/list")

  // 热门排行榜
  PopularResp Popular(1: PopularReq req) (api.get = "/video/popular")

  // 搜索
  SearchResp Search(1: SearchReq req) (api.post = "/video/search")
}

