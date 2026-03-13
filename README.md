# tiktok-go

@west2-online go组work4
一个基于[Hertz](https://github.com/cloudwego/hertz)框架构建的乞丐版石山TikTok

## 技术栈

- 框架：Hertz
- 数据库：PostgreSQL (GORM Gen)  
  "构式Gen，下辈子不用了，子查询都写不了，何意味？"
- 缓存：Redis
- 认证：JWT (Access Token + Refresh Token)
- 配置：Viper

## 项目结构
```
tiktok-go/
├── main.go                        # 入口点
├── router_gen.go
├── config/                        # 配置模块
├── idl/                           # 接口定义
├── biz/                           # 业务逻辑层
│   ├── handler/                   # HTTP处理器
│   ├── service/                   # 业务服务层
│   │   ├── user/                  # 用户
│   │   ├── video/                 # 视频
│   │   ├── interaction/           # 互动
│   │   ├── social/                # 社交
│   │   └── chat/                  # 聊天
│   │
│   ├── router/                    # 路由
│   ├── model/                     # 请求和响应数据模型
│   ├── pack/                      # 响应数据打包方法
│   ├── mw/                        # 中间件
│   │   └── auth/                  # JWT认证中间件
│   │
│   └── chat/                      # WebSocket聊天处理
│
├── pkg/                           # 通用工具包
│   ├── db/                        # 数据库访问
│   │   ├── postgres.go            # 数据库连接初始化
│   │   ├── model/                 # 数据库模型
│   │   ├── query/                 # Gen查询
│   │   └── ...
│   │
│   ├── cache/                     # 缓存
│   │   ├── redis.go               # Redis连接初始化
│   │   └── ...
│   │
│   ├── bucket/                    # 对象存储
│   │   ├── minio.go               # MinIO客户端初始化
│   │   └── ...
│   │
│   ├── jwt/                       # JWT工具包
│   ├── ffmpeg/                    # 媒体处理
│   ├── img/                       # 图片处理
│   ├── errno/                     # 错误码定义
│   ├── constants/                 # 全局常量
│   ├── totp/                      # TOTP多因素认证
│   └── utils/                     # 工具函数          
│
└── cmd/
    └── gorm-gen/                  # GORM Gen
```

## API文档

[文档](./API.md)

## 总结文档

[总结](./Summary.md)

## 快速开始

1. 复制配置文件并填写配置：
   ```bash
   cp config/config.yaml.example config/config.yaml
   ```

2. 启动服务：
   ```bash
   go run main.go
   ```
   
## Docker部署

1. 启动容器:
    ```bash
    docker run -d --name tiktok-go -p 13215:13215 -v C:/docker/tiktok-go/config:/app/config tiktok-go
    ```
   
2. 修改配置文件 `config/config.yaml`

3. 重启容器:
    ```bash
    docker restart tiktok-go
    ```

---

## TODO List

### 用户模块 (User)

- [x] `POST /user/register` — 用户注册
- [x] `POST /user/login` — 用户登录
- [x] `POST /auth/refresh` — 刷新 Token
- [X] `GET  /user/info` — 获取用户信息
- [X] `PUT  /user/avatar/upload` — 上传头像
- [x] `GET  /auth/mfa/qrcode` — 获取 MFA 二维码
- [x] `POST /auth/mfa/bind` — 绑定 MFA
- [ ] `POST /user/image/search` — 以图搜图 (没看懂，何意味)

### 视频模块 (Video)

- [X] `GET  /video/feed` — 视频流
- [X] `POST /video/publish` — 发布视频
- [X] `GET  /video/list` — 用户视频列表
- [X] `GET  /video/popular` — 热门视频
- [X] `POST /video/search` — 搜索视频
- [X] `GET  /video/visit` — 访问视频

### 互动模块 (Interaction)

- [X] `POST   /like/action` — 点赞 / 取消点赞
- [X] `GET    /like/list` — 点赞列表
- [X] `POST   /comment/publish` — 发布评论
- [X] `GET    /comment/list` — 评论列表
- [X] `DELETE /comment/delete` — 删除评论

### 社交模块 (Social)

- [X] `POST /relation/action` — 关注 / 取消关注
- [X] `GET  /following/list` — 关注列表
- [X] `GET  /follower/list` — 粉丝列表
- [x] `GET  /friends/list` — 好友列表
- [X] `GET  /ws` — WebSocket 聊天

