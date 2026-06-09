package flowlimit

import (
	"fmt"
	"simplest_script/core"
	"simplest_script/core/svc"
	"time"
)

type PeriodItem struct {
	Limit    int
	Duration int
}

type FlowLimit struct {
	Flag   string
	Period []PeriodItem
}

func NewFlowLimit() *FlowLimit {
	return &FlowLimit{}
}

func (l *FlowLimit) Make(flag string) *FlowLimit {
	l.Flag = flag
	return l
}

func (l *FlowLimit) SetPeriod(limit int, duration int) *FlowLimit {
	l.Period = append(l.Period, PeriodItem{
		Limit:    limit,
		Duration: duration,
	})

	return l
}

func (l *FlowLimit) Add(num int) *FlowLimit {
	for _, v := range l.Period {
		key := fmt.Sprintf("flow_limit:%s_%d", l.Flag, v.Duration)
		count, err := svc.NewRedis(core.RDSDefault).Incr(key).Result()

		if err == nil && count == int64(num) {
			svc.NewRedis(core.RDSDefault).Expire(key, time.Duration(v.Duration)*time.Second)
		}
	}

	return l
}

func (l *FlowLimit) Check() bool {
	for _, v := range l.Period {
		key := fmt.Sprintf("flow_limit:%s_%d", l.Flag, v.Duration)
		count, err := svc.NewRedis(core.RDSDefault).Get(key).Int()

		if err != nil || count > v.Limit {
			return false
		}
	}

	return true
}
