# tiktok-go

@west2-online go组work4
一个基于[Hertz](https://github.com/cloudwego/hertz)框架构建的乞丐版TikTok

## 技术栈

- 框架：Cloudwego Hertz
- 数据库：PostgreSQL(GORM Gen)
- 缓存：Redis
- 认证：JWT(Access Token + Refresh Token)
- 配置：Viper

## 快速开始

1. 复制配置文件并填写配置：
   ```bash
   cp config/config.yaml.example config/config.yaml
   ```

2. 启动服务：
   ```bash
   go run main.go
   ```

---

## TODO List

### 用户模块(User)

- [x] `POST /user/register` — 用户注册
- [x] `POST /user/login` — 用户登录
- [x] `POST /auth/refresh` — 刷新 Token
- [X] `GET  /user/info` — 获取用户信息
- [X] `PUT  /user/avatar/upload` — 上传头像
- [x] `GET  /auth/mfa/qrcode` — 获取 MFA 二维码
- [x] `POST /auth/mfa/bind` — 绑定 MFA
- [ ] `POST /user/image/search` — 以图搜图

### 视频模块(Video)

- [X] `GET  /video/feed` — 视频流
- [X] `POST /video/publish` — 发布视频
- [X] `GET  /video/list` — 用户视频列表
- [X] `GET  /video/popular` — 热门视频
- [X] `POST /video/search` — 搜索视频
- [X] `GET  /video/visit` — 访问视频

### 互动模块(Interaction)

- [X] `POST   /like/action` — 点赞 / 取消点赞
- [X] `GET    /like/list` — 点赞列表
- [ ] `POST   /comment/publish` — 发布评论
- [ ] `GET    /comment/list` — 评论列表
- [ ] `DELETE /comment/delete` — 删除评论

### 社交模块(Social)

- [ ] `POST /relation/action` — 关注 / 取消关注
- [ ] `GET  /following/list` — 关注列表
- [ ] `GET  /follower/list` — 粉丝列表
- [ ] `GET  /friends/list` — 好友列表
- [ ] `GET  /ws` — WebSocket 聊天

