package warning

import (
	"simplest_script/crontab"
	warningService "simplest_script/internal/services/warning"
)

type MediaDataWarning struct {
}

func NewMediaDataWarning() *MediaDataWarning {
	return &MediaDataWarning{}
}

func (l *MediaDataWarning) Handler(params string) *crontab.Result {
	for k := range warningService.WarningTitleMap {
		warningService.NewWarning(k).SetPeriod(1, 600).SendMessage()
	}

	return &crontab.Result{
		Status: 0,
		Data:   "success",
	}
}
