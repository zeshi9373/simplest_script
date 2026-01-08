package data

import (
	"encoding/json"
	"simplest_script/core/tool"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	AccountReportUrl = "https://ad.e.kuaishou.com/rest/openapi/v1/report/account_report"
)

type KuaishouReportData struct {
}

func NewKuaishouReportData() *KuaishouReportData {
	return &KuaishouReportData{}
}

func (h *KuaishouReportData) AccountReport(accessToken string, param *AccountReportDataRequest, list *Data) (code int, err error) {
	header := map[string]string{
		"Access-Token": accessToken,
	}

	data, _ := json.Marshal(param)
	response, err := tool.NewHttp(AccountReportUrl, 5*time.Second).Post(header, data)

	if err != nil {
		hlog.Error("KS查询报表数据请求错误1 error: " + err.Error() + " param: " + string(data) + " response: " + string(response))
		return 500, err
	}

	res := &AccountReportDataBase{}
	err = json.Unmarshal(response, res)

	if err != nil {
		hlog.Error("KS查询报表数据解析错误2 error: " + err.Error() + " param: " + string(data) + "  response: " + string(response))
		return 500, err
	}

	if res.Code != 0 {
		// warningService.NewMediaData().SetPeriod(1, 600).Add("快手时段数据", string(response))
		hlog.Error("KS查询报表数据请求错误3  param: " + string(data) + " response: " + string(response))
		return res.Code, err
	}

	result := &AccountReportDataResponse{}
	err = json.Unmarshal(response, result)
	if err != nil {
		hlog.Error("KS查询报表数据解析错误3 error: " + err.Error() + " param: " + string(data) + "  response: " + string(response))
		return 500, err
	}

	hlog.Info("账户id", param.AdvertiserID, " param: "+string(data)+" response: ", string(response))
	*list = result.Data

	return 0, nil
}
