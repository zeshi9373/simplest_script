package warning

import (
	"fmt"
	"simplest_script/core"
	"simplest_script/core/svc"
	"simplest_script/core/tool"
	"strconv"
	"time"
)

const (
	WarningKey   = "script_warning:%s.%d"
	WarningNonce = "script_warning_nonce:%s.%d"
)

type PeriodItem struct {
	Limit    int
	Duration int
}

type Warning struct {
	Flag   string
	Period []PeriodItem
}

func NewWarning(flag string) *Warning {
	return &Warning{
		Flag: flag,
	}
}

func (w *Warning) SetPeriod(period []PeriodItem) *Warning {
	w.Period = append(w.Period, period...)
	return w
}

func (w *Warning) Add(num int64) *Warning {
	if num > 0 {
		for _, v := range w.Period {
			key := fmt.Sprintf(WarningKey, w.Flag, v.Duration)
			count, err := svc.NewRedis(core.RDSDefault).IncrBy(key, num).Result()

			if err == nil && count == num {
				svc.NewRedis(core.RDSDefault).Expire(key, time.Duration(v.Duration)*time.Second)
				svc.NewRedis(core.RDSDefault).Set(fmt.Sprintf(WarningNonce, w.Flag, v.Duration), tool.Uuid(), time.Duration(v.Duration)*time.Second)
			}
		}
	}

	return w
}

func (w *Warning) GetNonce(duration int) string {
	var nonce string

	for _, v := range w.Period {
		nonce = svc.NewRedis(core.RDSDefault).Get(fmt.Sprintf(WarningNonce, w.Flag, v.Duration)).Val()

		if v.Duration == duration {
			return nonce
		}
	}

	return nonce
}

func (w *Warning) GetCount(duration int) int64 {
	var count int64

	for _, v := range w.Period {
		countStr, err := svc.NewRedis(core.RDSDefault).Get(fmt.Sprintf(WarningKey, w.Flag, v.Duration)).Result()

		if err != nil {
			return 0
		}

		count, _ = strconv.ParseInt(countStr, 10, 64)

		if v.Duration == duration {
			return count
		}
	}

	return count
}

func (w *Warning) Check() (bool, int) {
	for _, v := range w.Period {
		count, err := svc.NewRedis(core.RDSDefault).Get(fmt.Sprintf(WarningKey, w.Flag, v.Duration)).Result()

		if err != nil {
			return false, 0
		}

		if countInt, err := strconv.ParseInt(count, 10, 64); err != nil || countInt >= int64(v.Limit) {
			return false, v.Duration
		}
	}

	return true, 0
}
