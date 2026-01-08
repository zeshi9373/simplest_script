package script

import (
	"simplest_script/crontab"
)

type HandlerFunc interface {
	Handler(params string) *crontab.Result
}

var HandlerEntry = make(map[string]HandlerFunc)

func Exec(args ...string) *crontab.Result {
	fn := args[0]
	params := args[1]
	InitEntry()

	if fn == "" {
		return &crontab.Result{
			Status: 1,
			Data:   "请输入脚本名称 " + fn,
		}
	}

	if handler, ok := HandlerEntry[fn]; ok {
		if handler == nil {
			return &crontab.Result{
				Status: 1,
				Data:   "请输入脚本方法不存在 " + fn,
			}
		}

		return handler.Handler(params)
	}

	return &crontab.Result{
		Status: 1,
		Data:   "请输入正确的脚本名称 " + fn,
	}
}
