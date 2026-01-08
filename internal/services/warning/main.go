package warning

import (
	"encoding/json"
	"fmt"
	"os"
	"simplest_script/core"
	"simplest_script/core/logger"
	"simplest_script/core/svc"
	coreWarning "simplest_script/core/warning"
	"simplest_script/expand/feishu"
	"simplest_script/internal/consts"
	"strings"
	"time"
)

type Warning struct {
	Flag      string
	FeishuKey string
	Period    []coreWarning.PeriodItem
}

func NewWarning(flag string) *Warning {
	return &Warning{
		Flag: flag,
	}
}

func (w *Warning) SetFeishuKey(key string) *Warning {
	w.FeishuKey = key
	return w
}

func (w *Warning) SetPeriod(limit int, duration int) *Warning {
	w.Period = append(w.Period, coreWarning.PeriodItem{
		Limit:    limit,
		Duration: duration,
	})

	return w
}

func (w *Warning) Add(ty string, data ...string) {
	if len(data) > 0 {
		dataLen := 0

		for _, v := range data {
			if len(strings.Join(strings.Fields(v), "")) == 0 {
				continue
			}

			dataLen++
		}

		if dataLen == 0 {
			return
		}

		warn := coreWarning.NewWarning(w.Flag).SetPeriod(w.Period).Add(int64(dataLen))

		for _, v := range w.Period {
			nonce := warn.GetNonce(v.Duration)
			nonceKey := w.Flag + "_nonce_data:%s"
			key := fmt.Sprintf(nonceKey, nonce)
			cacheStr, err := svc.NewRedis(core.RDSDefault).Get(key).Result()

			if err != nil {
				cacheStrNew, _ := json.Marshal(data)
				svc.NewRedis(core.RDSDefault).Set(key, string(cacheStrNew), time.Duration(v.Duration)*time.Second)
			} else {
				var cacheData []string
				json.Unmarshal([]byte(cacheStr), &cacheData)
				var dataMsg []string
				var l int

				if len(cacheData) < 20 {
					l = 20 - len(cacheData)
				}

				for _, vv := range data {
					if len(strings.Join(strings.Fields(vv), "")) == 0 {
						continue
					}

					logger.NewLogger("cpaWarningLog").Info("[" + ty + "]" + vv)

					if l > 0 {
						if len(vv) > 200 {
							dataMsg = append(dataMsg, "["+ty+"]"+vv[0:200])
						} else {
							dataMsg = append(dataMsg, "["+ty+"]"+vv)
						}
					}

					l--
				}

				if len(dataMsg) > 0 {
					cacheData = append(cacheData, dataMsg...)
					cacheStrNew, _ := json.Marshal(cacheData)
					svc.NewRedis(core.RDSDefault).Set(key, string(cacheStrNew), time.Duration(v.Duration)*time.Second)
				}
			}
		}
	}
}

func (w *Warning) SendMessage() {
	coreWarn := coreWarning.NewWarning(w.Flag).SetPeriod(w.Period)
	status, duration := coreWarn.Check()

	if !status {
		nonce := coreWarn.GetNonce(duration)

		if len(nonce) > 0 {
			nonceKey := w.Flag + "_nonce_data:%s"
			key := fmt.Sprintf(nonceKey, nonce)
			cacheStr, err := svc.NewRedis(core.RDSDefault).Get(key).Result()

			if err == nil && len(cacheStr) > 0 {
				count := coreWarn.GetCount(duration)
				warningData := make([]string, 0)
				json.Unmarshal([]byte(cacheStr), &warningData)
				msg := fmt.Sprintf("%s\r\n环境:%s\r\n机器:%s\r\n错误数量:%d\r\n详细内容:%s\r\n时间:%s\r\n详细日志记录请服务器查看", "【"+WarningTitleMap[w.Flag]+"】", os.Getenv("SCRIPT_ENV"), os.Getenv("SCRIPT_PARTITION"), count, strings.Join(warningData, "\r\n"), time.Now().Format("2006-01-02 15:04:05"))

				if w.FeishuKey == "" {
					w.FeishuKey = consts.FeishuKeyMainWarning
				}

				feishu.NewSendMessage(w.FeishuKey).Send(feishu.MsgTypeText, msg)
			}
		}
	}
}
