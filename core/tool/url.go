package tool

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// 判断字符串是否已经 URL 编码
func IsURLEncoded(s string) bool {
	// 尝试解码
	decoded, err := url.QueryUnescape(s)
	if err != nil {
		return false
	}
	// 如果解码后的字符串与原始字符串不同，则说明是 URL 编码的
	return decoded != s
}

type URLComponents struct {
	Scheme string
	Host   string
	Path   string
	Query  string
}

func ParseURLEx(rawURL string) (URLComponents, error) {
	components := URLComponents{}
	u, err := url.Parse(rawURL)
	if err != nil {
		return components, err
	}
	// Scheme
	components.Scheme = u.Scheme
	// Host 和 Port
	hostname := u.Hostname()
	if hostname != "" {
		components.Host = hostname
	}
	// Path
	components.Path = u.Path
	// Query
	components.Query = u.RawQuery

	return components, nil
}

func FormatParas(params map[string]any) string {
	if len(params) == 0 {
		return ""
	}

	// 获取并排序 keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, key := range keys {
		value := params[key]

		// 处理值
		valueStr := formatValueSimple(value)
		pairs = append(pairs, key+"="+valueStr)
	}

	return strings.Join(pairs, "&")
}

func formatValueSimple(value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		// 数组类型
		strs := make([]string, len(v))
		for i, item := range v {
			strs[i] = fmt.Sprintf("%v", item)
		}
		sort.Strings(strs)
		return strings.Join(strs, ",")

	case []string:
		// 字符串数组
		strs := make([]string, len(v))
		copy(strs, v)
		sort.Strings(strs)
		return strings.Join(strs, ",")

	case map[string]interface{}:
		// map 类型
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var pairs []string
		for _, k := range keys {
			pairs = append(pairs, k+"="+fmt.Sprintf("%v", v[k]))
		}

		encoded := strings.Join(pairs, "&")
		decoded, _ := url.QueryUnescape(encoded)
		return decoded

	default:
		// 基本类型
		return fmt.Sprintf("%v", v)
	}
}
