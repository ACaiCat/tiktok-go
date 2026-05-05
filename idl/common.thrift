namespace go common


// 响应状态
struct Base {
    // 返回消息
    1: required string msg;
    // 状态码
    2: required i32 code;
}

