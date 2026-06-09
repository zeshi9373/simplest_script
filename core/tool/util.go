package tool

import (
	"strconv"
)

// ConvertStringIdsToInt 将字符串ID数组转换为整数ID数组
func ConvertStringIdsToInt(ids []string) []int {
	if len(ids) == 0 {
		return []int{}
	}

	result := make([]int, 0)
	for _, idStr := range ids {
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		if id <= 0 {
			continue
		}

		result = append(result, id)
	}

	return result
}

// ConvertInt64IdsToString 将整数int64的ID数组转换为字符串ID数组
func ConvertInt64IdsToString(ids []int64) []string {
	if len(ids) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(ids))
	for _, id := range ids {
		// 跳过无效ID（根据业务逻辑，0通常表示无效ID）
		if id <= 0 {
			continue
		}

		result = append(result, strconv.FormatInt(id, 10))
	}

	return result
}

// 限制字符串长度，支持中英文混合
func LimitStringLength(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	// 按rune计算，确保中文字符正确处理
	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}

	return string(runes[:maxLength])
}
