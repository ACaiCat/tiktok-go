# API文档

鉴权方式: `JWT`，在`Authorization`请求头携带`AccessToken`

# 用户

## POST 刷新 Token

`POST /auth/refresh`

刷新 `access-token` 和 `refresh-token`。

> Body 请求参数

```json
{
  "refresh_token": "{{refresh_token}}"
}
```

### 请求参数

| 名称            | 位置   | 类型     | 必选 | 说明       |
|---------------|------|--------|----|----------|
| refresh_token | body | string | 是  | 刷新 token |

> 返回示例

```json
{
  "base": {
    "msg": "OK",
    "code": 10000
  },
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

### 返回数据结构

状态码 200

| 名称               | 类型      | 说明 |
|------------------|---------|----|
| » base           | object  | -  |
| »» msg           | string  | -  |
| »» code          | integer | -  |
| » data           | object  | -  |
| »» access_token  | string  | -  |
| »» refresh_token | string  | -  |

---

## POST 注册

`POST /user/register`

> Body 请求参数

```yaml
username: Cai
password: "114514"
```

### 请求参数

| 名称       | 位置   | 类型     | 必选 | 说明 |
|----------|------|--------|----|----|
| username | body | string | 是  | 账号 |
| password | body | string | 是  | 密码 |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  }
}
```

### 返回数据结构

状态码 200

| 名称      | 类型            | 说明 |
|---------|---------------|----|
| » base  | [响应状态](#响应状态) | -  |
| »» code | integer       | -  |
| »» msg  | string        | -  |

---

## POST 登录

`POST /user/login`

> Body 请求参数

```yaml
username: Cai
password: "114514"
code: ""
```

### 请求参数

| 名称       | 位置   | 类型     | 必选 | 说明             |
|----------|------|--------|----|----------------|
| username | body | string | 是  | 账号             |
| password | body | string | 是  | 密码             |
| code     | body | string | 否  | 验证码 / MFA Code |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "id": "483675103821824000",
    "username": "west2-online",
    "avatar_url": "https://west2.online/avatar-example.jpeg",
    "created_at": "1970-01-01 08:00:00",
    "updated_at": "1970-01-01 08:00:00",
    "deleted_at": "1970-01-01 08:00:00"
  }
}
```

### 返回数据结构

状态码 200

| 名称            | 类型            | 说明           |
|---------------|---------------|--------------|
| » base        | [响应状态](#响应状态) | -            |
| »» code       | integer       | -            |
| »» msg        | string        | -            |
| » data        | object        | -            |
| »» id         | string        | 用户唯一标识符      |
| »» username   | string        | 用户名          |
| »» avatar_url | string        | 头像链接（本地或云存储） |
| »» created_at | string        | 创建时间         |

---

## GET 用户信息

`GET /user/info`

### 请求参数

| 名称      | 位置    | 类型     | 必选 | 说明    |
|---------|-------|--------|----|-------|
| user_id | query | string | 否  | 用户 ID |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "id": "483675103821824000",
    "username": "west2-online",
    "avatar_url": "https://west2.online/avatar-example.jpeg",
    "created_at": "1970-01-01 08:00:00"
  }
}
```

### 返回数据结构

状态码 200

| 名称            | 类型            | 说明           |
|---------------|---------------|--------------|
| » base        | [响应状态](#响应状态) | -            |
| »» code       | integer       | -            |
| »» msg        | string        | -            |
| » data        | object        | -            |
| »» id         | string        | 用户唯一标识符      |
| »» username   | string        | 用户名          |
| »» avatar_url | string        | 头像链接（本地或云存储） |
| »» created_at | string        | 创建时间         |

---

## PUT 上传头像

`PUT /user/avatar/upload`

对当前用户上传头像。

> Body 请求参数

```yaml
data: ""
```

### 请求参数

| 名称            | 位置     | 类型             | 必选 | 说明                |
|---------------|--------|----------------|----|-------------------|
| data          | body   | string(binary) | 否  | 图片原始数据，服务端需校验文件类型 |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "id": "483675103821824000",
    "username": "west2-online",
    "avatar_url": "https://west2.online/avatar-example.jpeg",
    "created_at": "1970-01-01 08:00:00"
  }
}
```

### 返回数据结构

状态码 200

| 名称            | 类型            | 说明            |
|---------------|---------------|---------------|
| » base        | [响应状态](#响应状态) | -             |
| » data        | [用户](#用户)     | -             |
| »» id         | string        | 用户唯一标识符       |
| »» username   | string        | 用户名           |
| »» password   | string        | 密码（bcrypt 加密） |
| »» avatar_url | string        | 头像链接（本地或云存储）  |
| »» created_at | string        | 创建时间          |

---

## GET 获取 MFA 二维码

`GET /auth/mfa/qrcode`

获取绑定 MFA 时所需的二维码。

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "Success"
  },
  "data": {
    "secret": "2dGJT1gjtzo4zNybLa9A",
    "qrcode": "data:image/png;base64,..."
  }
}
```

### 返回数据结构

状态码 200

| 名称        | 类型            | 说明              |
|-----------|---------------|-----------------|
| » base    | [响应状态](#响应状态) | -               |
| » data    | object        | -               |
| »» secret | string        | 多重身份验证密钥        |
| »» qrcode | string        | base64 编码的二维码图片 |

---

## POST 绑定多因素身份认证 (MFA)

`POST /auth/mfa/bind`

> Body 请求参数

```yaml
code: ""
secret: ""
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明  |
|---------------|--------|--------|----|-----|
| code          | body   | string | 否  | 校验码 |
| secret        | body   | string | 否  | 密钥  |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  }
}
```

### 返回数据结构

状态码 200

| 名称      | 类型            | 说明 |
|---------|---------------|----|
| » base  | [响应状态](#响应状态) | -  |
| »» code | integer       | -  |
| »» msg  | string        | -  |

---

## POST 以图搜图

`POST /user/image/search`

用户上传图片原始数据，返回搜到的图片 URL。

> Body 请求参数

```yaml
data: ""
```

### 请求参数

| 名称   | 位置   | 类型             | 必选 | 说明     |
|------|------|----------------|----|--------|
| data | body | string(binary) | 否  | 图片原始数据 |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": "https://ice-pomelo-1318531687.cos.ap-nanjing.myqcloud.com/app_cover.png"
}
```

### 返回数据结构

状态码 200

| 名称      | 类型            | 说明     |
|---------|---------------|--------|
| » base  | [响应状态](#响应状态) | -      |
| »» code | integer       | -      |
| »» msg  | string        | -      |
| » data  | string        | 图片 URL |

---

# 视频

## GET 访问视频

`GET /video/visit`

用户访问视频计数。

### 请求参数

| 名称       | 位置    | 类型     | 必选 | 说明 |
|----------|-------|--------|----|----|
| video_id | query | string | 否  | -  |

> 返回示例

```json
{
  "base": null,
  "data": "string"
}
```

### 返回数据结构

状态码 200

| 名称     | 类型     | 说明 |
|--------|--------|----|
| » base | any    | -  |
| » data | string | -  |

---

## GET 视频流

`GET /video/feed/`

获取首页视频流。

### 请求参数

| 名称          | 位置    | 类型     | 必选 | 说明                      |
|-------------|-------|--------|----|-------------------------|
| latest_time | query | string | 否  | 13 位时间戳，若存在则返回此时间之后的视频流 |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "483675103821824000",
        "user_id": "483675103821824000",
        "video_url": "https://west2.online/video-example.mp4",
        "cover_url": "https://west2.online/video-example.jpeg",
        "title": "west2-online",
        "description": "west2-online",
        "visit_count": 0,
        "like_count": 0,
        "comment_count": 0,
        "created_at": "1970-01-01 08:00:00"
      }
    ]
  }
}
```

### 返回数据结构

状态码 200

| 名称                | 类型            | 说明                      |
|-------------------|---------------|-------------------------|
| » base            | [响应状态](#响应状态) | -                       |
| » data            | object        | -                       |
| »» items          | [[视频](#视频)]   | -                       |
| »»» id            | string        | 视频唯一标识符（可选自增/雪花/UUID 等） |
| »»» user_id       | string        | 发表视频的用户唯一标识符            |
| »»» video_url     | string        | 视频文件链接                  |
| »»» cover_url     | string        | 封面链接                    |
| »»» title         | string        | 视频标题                    |
| »»» description   | string        | 视频描述                    |
| »»» visit_count   | integer       | 访问量                     |
| »»» like_count    | integer       | 点赞数量                    |
| »»» comment_count | integer       | 评论数量                    |
| »»» created_at    | string        | 创建时间                    |

---

## POST 投稿

`POST /video/publish`

使用 HTTP 单文件上传视频。

> Body 请求参数

```yaml
data: [ ]
title: west2-online
description: west2@online
```

### 请求参数

| 名称            | 位置     | 类型             | 必选 | 说明     |
|---------------|--------|----------------|----|--------|
| data          | body   | string(binary) | 否  | 视频原始数据 |
| title         | body   | string         | 否  | 视频标题   |
| description   | body   | string         | 否  | 描述     |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  }
}
```

### 返回数据结构

状态码 200

| 名称      | 类型            | 说明 |
|---------|---------------|----|
| » base  | [响应状态](#响应状态) | -  |
| »» code | integer       | -  |
| »» msg  | string        | -  |

---

## GET 发布列表

`GET /video/list`

根据 user_id 查看指定人的发布列表。

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明        |
|---------------|--------|---------|----|-----------|
| user_id       | query  | string  | 是  | -         |
| page_num      | query  | integer | 是  | 页码，从 0 开始 |
| page_size     | query  | integer | 是  | 单页尺寸      |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "483675103821824000",
        "user_id": "483675103821824000",
        "video_url": "https://west2.online/video-example.mp4",
        "cover_url": "https://west2.online/video-example.jpeg",
        "title": "west2-online",
        "description": "west2-online",
        "visit_count": 0,
        "like_count": 0,
        "comment_count": 0,
        "created_at": "1970-01-01 08:00:00"
      }
    ],
    "total": 100
  }
}
```

### 返回数据结构

状态码 200

| 名称                | 类型            | 说明                      |
|-------------------|---------------|-------------------------|
| » base            | [响应状态](#响应状态) | -                       |
| » data            | object        | -                       |
| »» items          | [[视频](#视频)]   | 视频列表                    |
| »»» id            | string        | 视频唯一标识符                 |
| »»» user_id       | string        | 发表视频的用户唯一标识符            |
| »»» video_url     | string        | 视频文件链接                  |
| »»» cover_url     | string        | 封面链接                    |
| »»» title         | string        | 视频标题                    |
| »»» description   | string        | 视频描述                    |
| »»» visit_count   | integer       | 访问量                     |
| »»» like_count    | integer       | 点赞数量                    |
| »»» comment_count | integer       | 评论数量                    |
| »»» created_at    | string        | 创建时间                    |
| »» total          | integer       | 当前用户的全部发布数量（非 items 数量） |

---

## GET 热门排行榜

`GET /video/popular`

根据点击量（visit_count）获取排行榜数据，使用 Redis 缓存。

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明   |
|---------------|--------|---------|----|------|
| page_size     | query  | integer | 否  | 单页尺寸 |
| page_num      | query  | integer | 否  | 页码   |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "483675103821824000",
        "user_id": "483675103821824000",
        "video_url": "https://west2.online/video-example.mp4",
        "cover_url": "https://west2.online/video-example.jpeg",
        "title": "west2-online",
        "description": "west2-online",
        "visit_count": 0,
        "like_count": 0,
        "comment_count": 0,
        "created_at": "1970-01-01 08:00:00"
      }
    ]
  }
}
```

### 返回数据结构

状态码 200

| 名称                | 类型            | 说明           |
|-------------------|---------------|--------------|
| » base            | [响应状态](#响应状态) | -            |
| » data            | object        | -            |
| »» items          | [[视频](#视频)]   | 视频列表         |
| »»» id            | string        | 视频唯一标识符      |
| »»» user_id       | string        | 发表视频的用户唯一标识符 |
| »»» video_url     | string        | 视频文件链接       |
| »»» cover_url     | string        | 封面链接         |
| »»» title         | string        | 视频标题         |
| »»» description   | string        | 视频描述         |
| »»» visit_count   | integer       | 访问量          |
| »»» like_count    | integer       | 点赞数量         |
| »»» comment_count | integer       | 评论数量         |
| »»» created_at    | string        | 创建时间         |

---

## POST 搜索视频

`POST /video/search`

搜索指定关键字的视频，从以下字段进行搜索：

- 标题（title）
- 描述（description）

> Body 请求参数

```yaml
keywords: west
page_size: 0
page_num: 0
from_date: 0
to_date: 0
username: ""
```

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明               |
|---------------|--------|---------|----|------------------|
| keywords      | body   | string  | 是  | 关键字，留空则不限制       |
| page_size     | body   | integer | 是  | 单页尺寸             |
| page_num      | body   | integer | 是  | 页码               |
| from_date     | body   | integer | 否  | 起始时间，13 位时间戳     |
| to_date       | body   | integer | 否  | 结束时间，13 位时间戳     |
| username      | body   | string  | 否  | 查询包含指定字符的用户发布的视频 |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "483675103821824000",
        "user_id": "483675103821824000",
        "video_url": "https://west2.online/video-example.mp4",
        "cover_url": "https://west2.online/video-example.jpeg",
        "title": "west2-online",
        "description": "west2-online",
        "visit_count": 0,
        "like_count": 0,
        "comment_count": 0,
        "created_at": "1970-01-01 08:00:00"
      }
    ],
    "total": 100
  }
}
```

### 返回数据结构

状态码 200

| 名称                | 类型            | 说明           |
|-------------------|---------------|--------------|
| » base            | [响应状态](#响应状态) | -            |
| » data            | object        | -            |
| »» items          | [[视频](#视频)]   | 视频列表         |
| »»» id            | string        | 视频唯一标识符      |
| »»» user_id       | string        | 发表视频的用户唯一标识符 |
| »»» video_url     | string        | 视频文件链接       |
| »»» cover_url     | string        | 封面链接         |
| »»» title         | string        | 视频标题         |
| »»» description   | string        | 视频描述         |
| »»» visit_count   | integer       | 访问量          |
| »»» like_count    | integer       | 点赞数量         |
| »»» comment_count | integer       | 评论数量         |
| »»» created_at    | string        | 创建时间         |
| »» total          | integer       | 查询结果总数量      |

---

# 互动

## POST 点赞操作

`POST /like/action`

> Body 请求参数

```yaml
video_id: "483675103821824000"
comment_id: "483675103821824000"
action_type: "1"
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明                             |
|---------------|--------|--------|----|--------------------------------|
| video_id      | body   | string | 否  | video_id 和 comment_id 必须存在其中一个 |
| comment_id    | body   | string | 否  | video_id 和 comment_id 必须存在其中一个 |
| action_type   | body   | string | 否  | 1：点赞，2：取消点赞                    |

#### 枚举值

| 属性          | 值 |
|-------------|---|
| action_type | 1 |
| action_type | 2 |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  }
}
```

```json
{
  "base": {
    "code": -1,
    "msg": "点赞失败"
  }
}
```

### 返回数据结构

状态码 200

| 名称      | 类型      | 说明 |
|---------|---------|----|
| » base  | object  | -  |
| »» code | integer | -  |
| »» msg  | string  | -  |

---

## GET 点赞列表

`GET /like/list`

返回指定用户点赞的视频。

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明   |
|---------------|--------|---------|----|------|
| user_id       | query  | string  | 否  | -    |
| page_size     | query  | integer | 否  | 单页尺寸 |
| page_num      | query  | integer | 否  | 页码   |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "483675103821824000",
        "user_id": "483675103821824000",
        "video_url": "https://west2.online/video-example.mp4",
        "cover_url": "https://west2.online/video-example.jpeg",
        "title": "west2-online",
        "description": "west2-online",
        "visit_count": 0,
        "like_count": 0,
        "comment_count": 0,
        "created_at": "1970-01-01 08:00:00"
      }
    ]
  }
}
```

### 返回数据结构

状态码 200

| 名称                | 类型          | 说明           |
|-------------------|-------------|--------------|
| » base            | object      | -            |
| »» code           | integer     | -            |
| »» msg            | string      | -            |
| » data            | object      | -            |
| »» items          | [[视频](#视频)] | 视频列表         |
| »»» id            | string      | 视频唯一标识符      |
| »»» user_id       | string      | 发表视频的用户唯一标识符 |
| »»» video_url     | string      | 视频文件链接       |
| »»» cover_url     | string      | 封面链接         |
| »»» title         | string      | 视频标题         |
| »»» description   | string      | 视频描述         |
| »»» visit_count   | integer     | 访问量          |
| »»» like_count    | integer     | 点赞数量         |
| »»» comment_count | integer     | 评论数量         |
| »»» created_at    | string      | 创建时间         |

---

## POST 评论

`POST /comment/publish`

对视频或评论进行评论。

> Body 请求参数

```yaml
video_id: "483675103821824000"
comment_id: "483675103821824000"
content: 我想睡觉
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明                             |
|---------------|--------|--------|----|--------------------------------|
| video_id      | body   | string | 否  | video_id 和 comment_id 必须存在其中一个 |
| comment_id    | body   | string | 否  | video_id 和 comment_id 必须存在其中一个 |
| content       | body   | string | 是  | 评论内容                           |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  }
}
```

### 返回数据结构

状态码 200

| 名称      | 类型      | 说明 |
|---------|---------|----|
| » base  | object  | -  |
| »» code | integer | -  |
| »» msg  | string  | -  |

---

## GET 评论列表

`GET /comment/list`

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明                             |
|---------------|--------|---------|----|--------------------------------|
| video_id      | query  | string  | 否  | video_id 和 comment_id 必须存在其中一个 |
| comment_id    | query  | string  | 否  | video_id 和 comment_id 必须存在其中一个 |
| page_size     | query  | integer | 否  | 单页尺寸                           |
| page_num      | query  | integer | 否  | 页码                             |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "820000201009032177",
        "user_id": "330000198910053439",
        "video_id": "370000201804261922",
        "parent_id": "210000202301182026",
        "like_count": 89,
        "child_count": 89,
        "content": "这视频真好看。",
        "created_at": "2020-07-26 23:32:44"
      }
    ]
  }
}
```

### 返回数据结构

状态码 200

| 名称              | 类型          | 说明             |
|-----------------|-------------|----------------|
| » base          | object      | -              |
| »» code         | integer     | -              |
| »» msg          | string      | -              |
| » data          | object      | -              |
| »» items        | [[评论](#评论)] | 评论列表           |
| »»» id          | string      | 评论唯一标识符        |
| »»» user_id     | string      | 发表评论的用户唯一标识符   |
| »»» video_id    | string      | 视频唯一标识符        |
| »»» parent_id   | string      | 父评论唯一标识符       |
| »»» like_count  | integer     | 点赞数量           |
| »»» child_count | integer     | 子评论数量          |
| »»» content     | string      | 评论内容（建议进行文本处理） |
| »»» created_at  | string      | 创建时间           |

---

## DELETE 删除评论

`DELETE /comment/delete`

> Body 请求参数

```yaml
comment_id: "1"
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明    |
|---------------|--------|--------|----|-------|
| comment_id    | body   | string | 否  | 评论 ID |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  }
}
```

### 返回数据结构

状态码 200

| 名称      | 类型      | 说明 |
|---------|---------|----|
| » base  | object  | -  |
| »» code | integer | -  |
| »» msg  | string  | -  |

---

# 社交

## POST 关注操作

`POST /relation/action`

> Body 请求参数

```yaml
to_user_id: "1"
action_type: 1
```

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明        |
|---------------|--------|---------|----|-----------|
| to_user_id    | body   | string  | 是  | 操作对象用户 ID |
| action_type   | body   | integer | 是  | 1：关注，2：取关 |

#### 枚举值

| 属性          | 值 |
|-------------|---|
| action_type | 1 |
| action_type | 2 |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  }
}
```

### 返回数据结构

状态码 200

| 名称      | 类型            | 说明 |
|---------|---------------|----|
| » base  | [响应状态](#响应状态) | -  |
| »» code | integer       | -  |
| »» msg  | string        | -  |

---

## GET 关注列表

`GET /following/list`

根据 user_id 查看指定人的关注列表。

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明        |
|---------------|--------|---------|----|-----------|
| user_id       | query  | string  | 是  | -         |
| page_num      | query  | integer | 否  | 页码，从 0 开始 |
| page_size     | query  | integer | 否  | 单页尺寸      |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "483675103821824000",
        "username": "west2-online",
        "avatar_url": "https://west2.online/avatar-example.jpeg"
      }
    ],
    "total": 16
  }
}
```

### 返回数据结构

状态码 200

| 名称             | 类型              | 说明                      |
|----------------|-----------------|-------------------------|
| » base         | [响应状态](#响应状态)   | -                       |
| »» code        | integer         | -                       |
| »» msg         | string          | -                       |
| » data         | object          | -                       |
| »» items       | [[社交对象](#社交对象)] | 关注列表                    |
| »»» id         | string          | 用户唯一标识符                 |
| »»» username   | string          | 用户名                     |
| »»» avatar_url | string          | 头像链接（本地或云存储）            |
| »» total       | integer         | 当前用户的全部关注数量（非 items 数量） |

---

## GET 粉丝列表

`GET /follower/list`

根据 user_id 查看指定人的粉丝列表。

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明        |
|---------------|--------|---------|----|-----------|
| user_id       | query  | string  | 是  | -         |
| page_num      | query  | integer | 否  | 页码，从 0 开始 |
| page_size     | query  | integer | 否  | 单页尺寸      |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "483675103821824000",
        "username": "west2-online",
        "avatar_url": "https://west2.online/avatar-example.jpeg"
      }
    ],
    "total": 16
  }
}
```

### 返回数据结构

状态码 200

| 名称             | 类型              | 说明                      |
|----------------|-----------------|-------------------------|
| » base         | [响应状态](#响应状态)   | -                       |
| »» code        | integer         | -                       |
| »» msg         | string          | -                       |
| » data         | object          | -                       |
| »» items       | [[社交对象](#社交对象)] | 粉丝列表                    |
| »»» id         | string          | 用户唯一标识符                 |
| »»» username   | string          | 用户名                     |
| »»» avatar_url | string          | 头像链接（本地或云存储）            |
| »» total       | integer         | 当前用户的全部粉丝数量（非 items 数量） |

---

## GET 好友列表

`GET /friends/list`

查看当前登录用户的好友列表。

### 请求参数

| 名称            | 位置     | 类型      | 必选 | 说明        |
|---------------|--------|---------|----|-----------|
| page_num      | query  | integer | 否  | 页码，从 0 开始 |
| page_size     | query  | integer | 否  | 单页尺寸      |

> 返回示例

```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "items": [
      {
        "id": "483675103821824000",
        "username": "west2-online",
        "avatar_url": "https://west2.online/avatar-example.jpeg"
      }
    ],
    "total": 16
  }
}
```

### 返回数据结构

状态码 200

| 名称             | 类型              | 说明                      |
|----------------|-----------------|-------------------------|
| » base         | [响应状态](#响应状态)   | -                       |
| »» code        | integer         | -                       |
| »» msg         | string          | -                       |
| » data         | object          | -                       |
| »» items       | [[社交对象](#社交对象)] | 好友列表                    |
| »»» id         | string          | 用户唯一标识符                 |
| »»» username   | string          | 用户名                     |
| »»» avatar_url | string          | 头像链接（本地或云存储）            |
| »» total       | integer         | 当前用户的全部好友数量（非 items 数量） |

---

# 数据模型

## 视频

```json
{
  "id": "string",
  "user_id": "string",
  "video_url": "string",
  "cover_url": "string",
  "title": "string",
  "description": "string",
  "visit_count": 0,
  "like_count": 0,
  "comment_count": 0,
  "created_at": "string"
}
```

### 属性

| 名称            | 类型      | 说明                      |
|---------------|---------|-------------------------|
| id            | string  | 视频唯一标识符（可选自增/雪花/UUID 等） |
| user_id       | string  | 发表视频的用户唯一标识符            |
| video_url     | string  | 视频文件链接                  |
| cover_url     | string  | 封面链接                    |
| title         | string  | 视频标题                    |
| description   | string  | 视频描述                    |
| visit_count   | integer | 访问量                     |
| like_count    | integer | 点赞数量                    |
| comment_count | integer | 评论数量                    |
| created_at    | string  | 创建时间                    |

---

## 用户

```json
{
  "id": "string",
  "username": "string",
  "password": "string",
  "avatar_url": "string",
  "created_at": "string"
}
```

### 属性

| 名称         | 类型     | 说明            |
|------------|--------|---------------|
| id         | string | 用户唯一标识符       |
| username   | string | 用户名           |
| password   | string | 密码（bcrypt 加密） |
| avatar_url | string | 头像链接（本地或云存储）  |
| created_at | string | 创建时间          |

---

## 评论

```json
{
  "id": "string",
  "user_id": "string",
  "video_id": "string",
  "parent_id": "string",
  "like_count": 0,
  "child_count": 0,
  "content": "string",
  "created_at": "string"
}
```

### 属性

| 名称          | 类型      | 说明                      |
|-------------|---------|-------------------------|
| id          | string  | 评论唯一标识符（可选自增/雪花/UUID 等） |
| user_id     | string  | 发表评论的用户唯一标识符            |
| video_id    | string  | 视频唯一标识符                 |
| parent_id   | string  | 父评论唯一标识符                |
| like_count  | integer | 点赞数量                    |
| child_count | integer | 子评论数量                   |
| content     | string  | 评论内容（建议进行文本处理）          |
| created_at  | string  | 创建时间                    |

---

## 社交对象

<a id="schema社交对象"></a>

```json
{
  "id": "string",
  "username": "string",
  "avatar_url": "string"
}
```

### 属性

| 名称         | 类型     | 说明           |
|------------|--------|--------------|
| id         | string | 用户唯一标识符      |
| username   | string | 用户名          |
| avatar_url | string | 头像链接（本地或云存储） |

---

## 响应状态

<a id="schema响应状态"></a>

```json
{
  "code": 0,
  "msg": "string"
}
```

### 属性

| 名称   | 类型      | 说明   |
|------|---------|------|
| code | integer | 状态码  |
| msg  | string  | 状态信息 |

---

## WebSocket 聊天接口

`WebSocket /ws`

### 鉴权

请求头 `Authorization` 中携带 Access Token：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 基本消息格式

```json
{
  "type": int,
  "body": {
    ...
  }
}
```

### 消息类型

#### [1] 发送聊天消息

##### C <-> S

| 字段          | 类型     | 说明            |
|-------------|--------|---------------|
| sender_id   | int64  | 发送方用户 ID      |
| receiver_id | int64  | 接收方用户 ID      |
| content     | string | 消息内容          |
| timestamp   | int64  | 发送时间戳（Unix 秒） |

```json
{
  "type": 1,
  "body": {
    "sender_id": 1,
    "receiver_id": 2,
    "content": "你好，在吗？",
    "timestamp": 1741737600
  }
}
```

---

#### [2] 拉取历史消息

##### C → S

| 字段        | 类型    | 说明         |
|-----------|-------|------------|
| user_id   | int64 | 对方用户 ID    |
| page      | int   | 页码（从 1 开始） |
| page_size | int   | 每页条数       |

```json
{
  "type": 2,
  "body": {
    "user_id": 2,
    "page": 1,
    "page_size": 20
  }
}
```

##### S → C

| 字段       | 类型            | 说明           |
|----------|---------------|--------------|
| messages | ChatMessage[] | 历史消息列表，按时间排序 |

```json
{
  "type": 2,
  "body": {
    "messages": [
      {
        "sender_id": 1,
        "receiver_id": 2,
        "content": "你好，在吗？",
        "timestamp": 1741737600
      },
      {
        "sender_id": 2,
        "receiver_id": 1,
        "content": "在的！",
        "timestamp": 1741737720
      }
    ]
  }
}
```

---

#### [3] 拉取未读消息

##### C → S

| 字段      | 类型    | 说明             |
|---------|-------|----------------|
| user_id | int64 | 对方用户 ID（消息来源方） |

```json
{
  "type": 3,
  "body": {
    "user_id": 2
  }
}
```

##### S → C

```json
{
  "type": 3,
  "body": {
    "messages": [
      {
        "sender_id": 2,
        "receiver_id": 1,
        "content": "有空吗？",
        "timestamp": 1741737800
      },
      {
        "sender_id": 2,
        "receiver_id": 1,
        "content": "一起吃饭？",
        "timestamp": 1741737860
      }
    ]
  }
}
```

---

#### [4] 错误处理

##### S → C

| 字段      | 类型     | 说明   |
|---------|--------|------|
| code    | int    | 错误码  |
| message | string | 错误信息 |

```json
{
  "type": 4,
  "body": {
    "code": 17002,
    "message": "你们还不是好友"
  }
}
```
