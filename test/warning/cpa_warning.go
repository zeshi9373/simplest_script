package warning

import (
	warningService "simplest_script/internal/services/warning"
)

func Warning() {
	warningService.NewWarning(warningService.TestWarningKey).SetPeriod(1, 600).Add("测试媒体类型", "媒体回传出错")
}

func SendWarning() {
	for k := range warningService.WarningTitleMap {
		warningService.NewWarning(k).SetPeriod(1, 600).SendMessage()
	}
}
