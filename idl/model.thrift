namespace go model

// 用户
struct User {
    // ID
    1: required string id;
    // 用户名
    2: required string username;
    // 头像URL
    3: required string avatar_url;
    // 创建时间
    4: required string created_at;
}

// 视频
struct Video {
    // ID
    1: required string id;
    // 发布者用户ID
    2: required string user_id;
    // 视频URL
    3: required string video_url;
    // 封面URL
    4: required string cover_url;
    // 标题
    5: required string title;
    // 简介
    6: required string description;
    // 播放量
    7: required i64 visit_count;
    // 点赞数
    8: required i64 like_count;
    // 评论数
    9: required i64 comment_count;
    // 发布时间
    10: required string created_at;
}

// 视频列表
struct VideoList {
    // 视频列表
    1: required list<Video> items;
}

// 带总数的视频列表
struct VideoListWithTotal {
    // 视频列表
    1: required list<Video> items;
    2: required i64 total;
}

// 评论
struct Comment {
    // ID
    1: required string id;
    // 发布者用户ID
    2: required string user_id;
    // 所属视频ID
    3: required string video_id;
    // 所属评论ID
    4: required string parent_id;
    // 点赞数
    5: required i64 like_count;
    // 子评论数
    6: required i64 child_count;
    // 内容
    7: required string content;
    // 发布时间
    8: required string created_at;
}

// 评论数据
struct CommentList {
    // 评论列表
    1: required list<Comment> items;
}

// 社交用户
struct SocialUser {
    // ID
    1: required string id;
    // 用户名
    2: required string username;
    // 头像URL
    3: required string avatar_url;
}

// 社交用户数据
struct SocialUserListWithTotal {
    // 社交用户列表
    1: required list<SocialUser> items;
    // 总数
    2: required i32 total;
}

// 关注操作类型
enum FollowActionType {
    // 关注
    FOLLOW = 1,
    // 取消关注
    UNFOLLOW = 2
}

enum LikeActionType {
    ADD = 1,
    DELETE = 2
}