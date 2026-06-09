package tool

func StringsLimitLength(str string, length int) string {
	if len(str) > length {
		return str[0:length]
	}
	return str
}
