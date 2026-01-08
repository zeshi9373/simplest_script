package tool

func IsInIntSlice(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}

	return false
}

// IsInSlice 检查值是否在切片中，支持任意可比较类型
func IsInSlice[T comparable](slice []T, val T) bool {
	if len(slice) == 0 {
		return false
	}

	for _, item := range slice {
		if item == val {
			return true
		}
	}

	return false
}
