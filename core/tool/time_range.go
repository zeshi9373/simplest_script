package tool

import (
	"encoding/json"
	"fmt"
	"time"
)

const DateLayout = "2006-01-02"
const DatetimeLayout = "2006-01-02 15:04:05"
const TIMEZONE = "Asia/Shanghai" // 时区：
type TimeRangeParams struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// ParseTimeRange 解析参数并返回 Unix 时间戳（秒）的开始和结束时间
// 如果 params 为空或无效，则使用默认时间范围（最近 defaultHoursAgo 小时对应的自然日）
// 返回 startTime, endTime, error
func ParseTimeRange(params string, defaultHoursAgo int) (int64, int64, error) {

	if defaultHoursAgo < 0 {
		defaultHoursAgo = 72
	}

	var paramsData TimeRangeParams

	if params != "" {
		if err := json.Unmarshal([]byte(params), &paramsData); err != nil {
			return 0, 0, fmt.Errorf("failed to parse params: %w", err)
		}
	}

	// 加载时区
	timeLocation, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to load timezone: %w", err)
	}

	now := time.Now().In(timeLocation) // 使用指定时区
	defaultStart := now.Add(-time.Duration(defaultHoursAgo) * time.Hour).Format(DateLayout)
	defaultEnd := now.Format(DateLayout)

	startDate := defaultStart
	endDate := defaultEnd

	// 校验 StartDate
	if paramsData.StartDate != "" {
		if _, err := time.Parse(DateLayout, paramsData.StartDate); err != nil {
			return 0, 0, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = paramsData.StartDate
	}

	// 校验 EndDate
	if paramsData.EndDate != "" {
		if _, err := time.Parse(DateLayout, paramsData.EndDate); err != nil {
			return 0, 0, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD")
		}
		endDate = paramsData.EndDate
	}

	// 构造完整时间字符串
	startTimeStr := startDate + " 00:00:00"
	endTimeStr := endDate + " 23:59:59"

	// 使用指定时区解析时间
	startTimeParse, err := time.ParseInLocation(DatetimeLayout, startTimeStr, timeLocation)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse start time: %w", err)
	}

	endTimeParse, err := time.ParseInLocation(DatetimeLayout, endTimeStr, timeLocation)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse end time: %w", err)
	}

	return startTimeParse.Unix(), endTimeParse.Unix(), nil
}

// 把时间字符串类型 2022-05-16 10:00:01 转 time.Time  如 2022-05-16 10:00:01
func DateStringToTime(timeStr string) time.Time {
	if len(timeStr) == 0 {
		// 默认当前的时间
		timeStr = time.Now().Format(DatetimeLayout)
	}

	TimeLocation, _ := time.LoadLocation(TIMEZONE) //获取北京时间时区，很重要！

	tt, _ := time.ParseInLocation(DatetimeLayout, timeStr, TimeLocation) //
	//  tt.Unix() 时间戳
	//时间戳 to 时间
	tm := time.Unix(tt.Unix(), 0)
	return tm
}
