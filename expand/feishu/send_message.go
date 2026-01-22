package feishu

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"simplest_script/core"
	"simplest_script/core/tool"
	"time"
)

const (
	MsgTypeText = "text"

	SendMessageUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/"
)

type SendMessage struct {
	Key string
}

func NewSendMessage(key string) *SendMessage {
	return &SendMessage{
		Key: key,
	}
}

func (l *SendMessage) Send(msgType string, msg string) error {
	if msgType == "" || msg == "" || l.Key == "" {
		return errors.New("请检查消息类型，内容和Key")
	}

	if os.Getenv("SCRIPT_ENV") != core.EnvRelease {
		fmt.Println("飞书消息： ", msg)
		return nil
	}

	data := map[string]interface{}{
		"msg_type": msgType,
		"content":  map[string]string{"text": msg},
	}

	str, _ := json.Marshal(data)

	_, err := tool.NewHttp(SendMessageUrl+l.Key, 5*time.Second).Post(nil, str)

	return err
}
