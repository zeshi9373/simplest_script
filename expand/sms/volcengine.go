package sms

import (
	"github.com/volcengine/volc-sdk-golang/service/sms"
)

var testAk = "ak"
var testSk = "sk"

type VolcEngine struct {
}

func NewVolcEngine() *VolcEngine {
	return &VolcEngine{}
}

type VolcEngineSms struct {
}

func (v *VolcEngineSms) Send(phoneNumber, templateCode, templateParam string) (string, error) {
	sms.DefaultInstance.Client.SetAccessKey(testAk)
	sms.DefaultInstance.Client.SetSecretKey(testSk)
	req := &sms.SmsRequest{
		SmsAccount:    "smsAccount",
		Sign:          "sign",
		TemplateID:    "ST_xxx",
		TemplateParam: "",
		PhoneNumbers:  "188xxxxxxxx",
		Tag:           "tag",
	}
	result, statusCode, err := sms.DefaultInstance.Send(req)
	if err != nil {
		println("sms send err:%s", err.Error())
		return "", err
	}
	println("result is :%s, statusCode is:%d", result, statusCode)
	return "", nil
}
