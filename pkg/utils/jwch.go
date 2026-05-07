package utils

import (
	"net/http"
	"strings"
)

func ParseCookies(rawData string) []*http.Cookie {
	var cookies []*http.Cookie
	maxSplitNumber := 2

	// 按照分号分割每个 Cookie
	pairs := strings.Split(rawData, ";")
	for _, pair := range pairs {
		// 去除空格并检查是否为空
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		// 按等号分割键和值
		parts := strings.SplitN(pair, "=", maxSplitNumber)
		if len(parts) != maxSplitNumber {
			continue // 如果格式不正确，跳过
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 创建 http.Cookie 并添加到切片
		cookie := &http.Cookie{
			Name:  key,
			Value: value,
		}
		cookies = append(cookies, cookie)
	}

	return cookies
}

func ParseCookiesToString(cookies []*http.Cookie) string {
	var cookieStrings []string
	for _, cookie := range cookies {
		cookieStrings = append(cookieStrings, cookie.Name+"="+cookie.Value)
	}
	return strings.Join(cookieStrings, "; ")
}
