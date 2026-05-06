CREATE TABLE IF NOT EXISTS users
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      TEXT        NOT NULL,
    password      TEXT        NOT NULL,
    avatar_url    TEXT,
    totp_secret   TEXT,
    jwch_id       TEXT,
    jwch_password TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_users_username_active
    ON users (username)
    WHERE deleted_at IS NULL;

COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '用户ID';
COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.password IS '密码';
COMMENT ON COLUMN users.avatar_url IS '头像URL';
COMMENT ON COLUMN users.totp_secret IS 'TOTP密钥';
COMMENT ON COLUMN users.jwch_id IS '教务处学号';
COMMENT ON COLUMN users.jwch_password IS '教务处密码';
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
    deleted_at    TIMESTAMPTZ,

    CONSTRAINT fk_videos_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT chk_videos_visit_count_nonnegative CHECK (visit_count >= 0),
    CONSTRAINT chk_videos_like_count_nonnegative CHECK (like_count >= 0),
    CONSTRAINT chk_videos_comment_count_nonnegative CHECK (comment_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_videos_user_id
    ON videos (user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_videos_created_at
    ON videos (created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_videos_user_id_created_at
    ON videos (user_id, created_at DESC)
    WHERE deleted_at IS NULL;

COMMENT ON TABLE videos IS '视频表';
COMMENT ON COLUMN videos.id IS '视频ID';
COMMENT ON COLUMN videos.user_id IS '用户ID';
COMMENT ON COLUMN videos.video_url IS '视频URL';
COMMENT ON COLUMN videos.cover_url IS '封面URL';
COMMENT ON COLUMN videos.title IS '视频标题';
COMMENT ON COLUMN videos.description IS '视频描述';
COMMENT ON COLUMN videos.visit_count IS '访问量';
COMMENT ON COLUMN videos.like_count IS '点赞数';
COMMENT ON COLUMN videos.comment_count IS '评论数';
COMMENT ON COLUMN videos.created_at IS '创建时间';
COMMENT ON COLUMN videos.updated_at IS '更新时间';
COMMENT ON COLUMN videos.deleted_at IS '删除时间';

CREATE TABLE IF NOT EXISTS comments
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id       BIGINT      NOT NULL,
    video_id      BIGINT      NOT NULL,
    parent_id     BIGINT,
    content       TEXT        NOT NULL,
    like_count    BIGINT      NOT NULL DEFAULT 0,
    comment_count BIGINT      NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_comments_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT fk_comments_video_id FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE RESTRICT,
    CONSTRAINT fk_comments_parent_id FOREIGN KEY (parent_id) REFERENCES comments (id) ON DELETE CASCADE,
    CONSTRAINT chk_comments_not_self_parent CHECK (parent_id IS NULL OR parent_id <> id),
    CONSTRAINT chk_comments_like_count_non_negative CHECK (like_count >= 0),
    CONSTRAINT chk_comments_comment_count_non_negative CHECK (comment_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_comments_video_top_created_at
    ON comments (video_id, created_at DESC)
    WHERE parent_id IS NULL;


CREATE INDEX IF NOT EXISTS idx_comments_parent_created_at
    ON comments (parent_id, created_at ASC)
    WHERE parent_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_comments_user_created_at
    ON comments (user_id, created_at DESC);

COMMENT ON TABLE comments IS '评论表';
COMMENT ON COLUMN comments.id IS '评论ID';
COMMENT ON COLUMN comments.user_id IS '用户ID';
COMMENT ON COLUMN comments.video_id IS '视频ID';
COMMENT ON COLUMN comments.parent_id IS '父评论ID';
COMMENT ON COLUMN comments.content IS '评论内容';
COMMENT ON COLUMN comments.like_count IS '点赞数';
COMMENT ON COLUMN comments.comment_count IS '评论数';
COMMENT ON COLUMN comments.created_at IS '创建时间';

CREATE TABLE IF NOT EXISTS likes
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    video_id   BIGINT,
    comment_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_likes_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT fk_likes_video_id FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE RESTRICT,
    CONSTRAINT fk_likes_comment_id FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
    CONSTRAINT chk_likes_target_exactly_one CHECK (
        (video_id IS NOT NULL AND comment_id IS NULL)
            OR
        (video_id IS NULL AND comment_id IS NOT NULL)
        )
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_likes_user_video
    ON likes (user_id, video_id)
    WHERE video_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_likes_user_comment
    ON likes (user_id, comment_id)
    WHERE comment_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_likes_video_created_at
    ON likes (video_id, created_at DESC)
    WHERE video_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_likes_comment_created_at
    ON likes (comment_id, created_at DESC)
    WHERE comment_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_likes_user_created_at
    ON likes (user_id, created_at DESC);

COMMENT ON TABLE likes IS '点赞表';
COMMENT ON COLUMN likes.id IS '点赞ID';
COMMENT ON COLUMN likes.user_id IS '用户ID';
COMMENT ON COLUMN likes.video_id IS '视频ID';
COMMENT ON COLUMN likes.comment_id IS '评论ID';
COMMENT ON COLUMN likes.created_at IS '创建时间';

CREATE TABLE IF NOT EXISTS followers
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    follower_id BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_followers_user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,

    CONSTRAINT fk_followers_follower_id
        FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE RESTRICT,

    CONSTRAINT chk_followers_not_self
        CHECK (user_id <> follower_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_followers_user_follower
    ON followers (user_id, follower_id);

CREATE INDEX IF NOT EXISTS idx_followers_user_created_at
    ON followers (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_followers_follower_created_at
    ON followers (follower_id, created_at DESC);

COMMENT ON TABLE followers IS '关注表';
COMMENT ON COLUMN followers.id IS '关注ID';
COMMENT ON COLUMN followers.user_id IS '用户ID';
COMMENT ON COLUMN followers.follower_id IS '粉丝ID';
COMMENT ON COLUMN followers.created_at IS '创建时间';

CREATE TABLE IF NOT EXISTS chat_messages
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sender_id   BIGINT      NOT NULL,
    receiver_id BIGINT      NOT NULL,
    content     TEXT        NOT NULL,
    read_at     TIMESTAMPTZ,
    is_ai       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_chat_messages_sender_id
        FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE RESTRICT,

    CONSTRAINT fk_chat_messages_receiver_id
        FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_chat_messages_receiver_unread_created_at
    ON chat_messages (receiver_id, created_at DESC)
    WHERE read_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_chat_messages_sender_created_at
    ON chat_messages (sender_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_chat_messages_receiver_created_at
    ON chat_messages (receiver_id, created_at DESC);

COMMENT ON TABLE chat_messages IS '聊天消息表';
COMMENT ON COLUMN chat_messages.id IS '消息ID';
COMMENT ON COLUMN chat_messages.sender_id IS '发送者ID';
COMMENT ON COLUMN chat_messages.receiver_id IS '接收者ID';
COMMENT ON COLUMN chat_messages.content IS '消息内容';
COMMENT ON COLUMN chat_messages.read_at IS '阅读时间';
COMMENT ON COLUMN chat_messages.is_ai IS '是否为AI的消息';
COMMENT ON COLUMN chat_messages.created_at IS '创建时间';