# tiktok-go

@west2-online go组work5
一个基于[Hertz](https://github.com/cloudwego/hertz)框架构建的乞丐版石山TikTok

## 技术栈

- 框架：Hertz
- 数据库：PostgreSQL (GORM Gen)
- 缓存：Redis
- 认证：JWT (双Token)
- 存储: Minio
- 配置：Viper

## 项目结构

```
tiktok-go/
├── cmd/
│   ├── api/                       # APP入口点
│   ├── chat/                      # 聊天客户端 (100 %AI)
│   └── gorm-gen/                  # GORM Gen
│
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
│   └── mw/                        # 中间件、
│       ├── log/                   # 日志中间件
│       └── auth/                  # JWT认证中间件
│   
└── pkg/                           # 通用工具包
    ├── db/                        # 数据库访问
    │   ├── postgres.go            # 数据库连接初始化
    │   ├── model/                 # 数据库模型
    │   ├── query/                 # Gen查询
    │   └── ...
    │
    ├── cache/                     # 缓存
    │   ├── redis.go               # Redis连接初始化
    │   └── ...
    │
    ├── bucket/                    # 对象存储
    │   ├── minio.go               # MinIO客户端初始化
    │   └── ...
    │
    ├── ai/                        # AI聊天
    ├── jwt/                       # JWT工具包
    ├── ffmpeg/                    # 媒体处理
    ├── img/                       # 图片处理
    ├── errno/                     # 错误码定义
    ├── constants/                 # 全局常量
    ├── totp/                      # TOTP多因素认证
    └── utils/                     # 工具函数          
```

## API文档

[文档](docs/API.md)
> ApiFox生成的

## 总结文档

[总结](docs/Report.md)
> Work4的，Work5的还没写

## 设计文档

[总结](docs/Design.md)
> 石山随笔

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

