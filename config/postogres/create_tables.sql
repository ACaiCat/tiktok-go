CREATE TABLE IF NOT EXISTS users
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username   TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '用户ID';
COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.password IS '密码';
COMMENT ON COLUMN users.avatar_url IS '头像URL';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';
COMMENT ON COLUMN users.deleted_at IS '删除时间';

CREATE TABLE IF NOT EXISTS videos
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id       BIGINT      NOT NULL,
    video_url     TEXT        NOT NULL,
    cover_url     TEXT        NOT NULL,
    title         TEXT        NOT NULL,
    description   TEXT        NOT NULL,
    visit_count   BIGINT      NOT NULL DEFAULT 0,
    like_count    BIGINT      NOT NULL DEFAULT 0,
    comment_count BIGINT      NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ
);

COMMENT ON TABLE videos IS '视频表';
COMMENT ON COLUMN videos.id IS '视频ID';
COMMENT ON COLUMN videos.user_id IS '用户ID';
COMMENT ON COLUMN videos.video_url IS '视频URL';
COMMENT ON COLUMN videos.cover_url IS '封面URL';
COMMENT ON COLUMN videos.title IS '视频标题';
COMMENT ON COLUMN videos.description IS '视频描述';
COMMENT ON COLUMN videos.visit_count IS '访问量';
COMMENT ON COLUMN videos.created_at IS '创建时间';
COMMENT ON COLUMN videos.updated_at IS '更新时间';
COMMENT ON COLUMN videos.deleted_at IS '删除时间';

CREATE TABLE IF NOT EXISTS comments
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    video_id   BIGINT      NOT NULL,
    parent_id  BIGINT,
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

COMMENT ON TABLE comments IS '评论表';
COMMENT ON COLUMN comments.id IS '评论ID';
COMMENT ON COLUMN comments.user_id IS '用户ID';
COMMENT ON COLUMN comments.video_id IS '视频ID';
COMMENT ON COLUMN comments.parent_id IS '父评论ID';
COMMENT ON COLUMN comments.content IS '评论内容';
COMMENT ON COLUMN comments.created_at IS '创建时间';
COMMENT ON COLUMN comments.updated_at IS '更新时间';
COMMENT ON COLUMN comments.deleted_at IS '删除时间';

CREATE TABLE IF NOT EXISTS likes
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    video_id   BIGINT      NOT NULL,
    comment_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

COMMENT ON TABLE likes IS '点赞表';
COMMENT ON COLUMN likes.id IS '点赞ID';
COMMENT ON COLUMN likes.user_id IS '用户ID';
COMMENT ON COLUMN likes.video_id IS '视频ID';
COMMENT ON COLUMN likes.comment_id IS '评论ID';
COMMENT ON COLUMN likes.created_at IS '创建时间';
COMMENT ON COLUMN likes.updated_at IS '更新时间';
COMMENT ON COLUMN likes.deleted_at IS '删除时间';