package data

import (
	"encoding/json"
	"simplest_script/core/tool"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	AccountReportUrl = "https://api.oceanengine.com/open_api/v3.0/report/custom/get/"
)

type OceanReportData struct {
}

func NewOceanReportData() *OceanReportData {
	return &OceanReportData{}
}

func (h *OceanReportData) AccountReport(accessToken string, param *AccountReportDataRequest, list *Data) error {
	header := map[string]string{
		"Access-Token": accessToken,
	}

	dimensions, _ := json.Marshal(param.Dimensions)
	metrics, _ := json.Marshal(param.Metrics)
	filters, _ := json.Marshal(param.Filters)
	orderBy, _ := json.Marshal(param.OrderBy)
	p := map[string]string{
		"dimensions":    string(dimensions),
		"advertiser_id": param.AdvertiserID,
		"metrics":       string(metrics),
		"filters":       string(filters),
		"page":          param.Page,
		"page_size":     param.PageSize,
		"start_time":    param.StartTime,
		"end_time":      param.EndTime,
		"order_by":      string(orderBy),
	}

	response, err := tool.NewHttp(AccountReportUrl, 5*time.Second).Get(header, p)

	if err != nil {
		hlog.Error("OCEAN查询报表数据请求错误1 error: " + err.Error() + " response: " + string(response))
		return err
	}

	res := &AccountReportDataBase{}
	err = json.Unmarshal(response, res)

	if err != nil {
		hlog.Error("OCEAN查询报表数据解析错误2 error: " + err.Error() + " response: " + string(response))
		return err
	}

	if res.Code == 40110 {
		time.Sleep(time.Millisecond * time.Duration(tool.Random(1000, 5000)))
		h.AccountReport(accessToken, param, list)
	}

	if res.Code != 0 {
		// warningService.NewMediaData().SetPeriod(1, 600).Add("快手时段数据", string(response))
		hlog.Error("OCEAN查询报表数据解析错误3 response: " + string(response))
		return err
	}

	result := &AccountReportDataResponse{}
	err = json.Unmarshal(response, result)

	if err != nil || res.Code != 0 {
		hlog.Error("OCEAN查询报表数据解析错误4 error: " + err.Error() + " response: " + string(response))
		return err
	}

	*list = result.Data

	return nil
}
