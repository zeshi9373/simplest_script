package tool

import (
	"math/rand"
	"time"
)

func Random(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandString(s []string, l int) string {
	if len(s) == 0 {
		return ""
	}

	if l <= 0 {
		return s[rand.Intn(len(s))]
	}

	var result string

	for i := 0; i < l; i++ {
		result += s[rand.Intn(len(s))]
	}

	return result
}

func GetTodayZeroTime() int64 {
	// 1. 加载上海时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Now().Truncate(24 * time.Hour).Unix()
	}

	// 2. 获取当前上海时间
	now := time.Now().In(loc)

	// 3. 构造当天零点（注意使用 loc 确保时区一致）
	zero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// 4. 获取 Unix 时间戳
	return zero.Unix()
}
