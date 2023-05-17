package utils

import (
	"math/rand"
	"time"

	"github.com/go-basic/uuid"
)

// 生成主键
func GetUuid() string {
	uuid := uuid.New()
	return uuid
}

//时间转时间戳
func Strtime2Int(datetime string) int64 {
	//日期转化为时间戳
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	tmp, _ := time.ParseInLocation(timeLayout, datetime, time.Local)
	timestamp := tmp.Unix() //转化为时间戳 类型是int64
	return timestamp
}

// 随机字符串
func GetRandomString(l int) string {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())
	// 生成随机字符串
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, l)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	// 输出结果
	return string(result)
}
