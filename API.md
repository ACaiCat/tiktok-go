# API文档

## 聊天接口 (/ws)

### 鉴权

请求头`Authorization`中携带Access Token

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

#### [1]发送聊天消息

##### C <-> S

| 字段            | 类型     | 说明            |
|---------------|--------|---------------|
| `sender_id`   | int64  | 发送方用户 ID      |
| `receiver_id` | int64  | 接收方用户 ID      |
| `content`     | string | 消息内容          |
| `timestamp`   | int64  | 发送时间戳（Unix 秒） |

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

#### [2]拉取历史消息

##### C → S

| 字段          | 类型    | 说明         |
|-------------|-------|------------|
| `user_id`   | int64 | 对方用户 ID    |
| `page`      | int   | 页码（从 1 开始） |
| `page_size` | int   | 每页条数       |

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

| 字段         | 类型            | 说明           |
|------------|---------------|--------------|
| `messages` | ChatMessage[] | 历史消息列表，按时间排序 |

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

#### [3]拉取未读消息

##### C → S

`body` 对应 `UnreadRequest`：

| 字段        | 类型    | 说明             |
|-----------|-------|----------------|
| `user_id` | int64 | 对方用户 ID（消息来源方） |

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

#### [4]错误处理

##### S → C

| 字段        | 类型     | 说明   |
|-----------|--------|------|
| `code`    | int    | 错误码  |
| `message` | string | 错误信息 |

```json
{
  "type": 4,
  "body": {
    "code": 17002,
    "message": "你们还不是好友"
  }
}
```

