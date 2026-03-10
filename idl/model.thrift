namespace go model

struct User {
  1: required string id;
  2: required string username;
  3: required string avatar_url;
  4: required string created_at;
}

struct Video {
  1: required string id;
  2: required string user_id;
  3: required string video_url;
  4: required string cover_url;
  5: required string title;
  6: required string description;
  7: required i64 visit_count;
  8: required i64 comment_count;
  9: required string created_at;
}

struct Comment {
  1: required string id;
  2: required string user_id;
  3: required string video_url;
  4: required string parent_id;
  5: required i64 like_count;
  6: required i64 child_count;
  7: required string content;
  8: required string created_at;
}

struct SocialUser {
  1: required string id;
  2: required string username;
  3: required string avatar_url;
}