package data

import (
	"encoding/json"
	"fmt"
	"simplest_script/core/tool"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	AccountReportUrl = "https://api.e.qq.com/v3.0/hourly_reports/get"
)

type GdtReportData struct {
}

func NewGdtReportData() *GdtReportData {
	return &GdtReportData{}
}

func (h *GdtReportData) AccountReport(accessToken string, param *AccountReportDataRequest, list *Data) error {
	groupBy, _ := json.Marshal(param.GroupBy)
	fields, _ := json.Marshal(param.Fields)
	p := map[string]string{
		"access_token": accessToken,
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
		"nonce":        tool.Md5(strconv.FormatInt(time.Now().UnixMicro(), 10) + tool.Uuid()),
		"account_id":   param.AccountID,
		"level":        param.Level,
		"page_size":    strconv.Itoa(param.PageSize),
		"page":         strconv.Itoa(param.Page),
		"date_range":   fmt.Sprintf("{\"start_date\":\"%s\",\"end_date\":\"%s\"}", param.DateRange.StartDate, param.DateRange.EndDate),
		"fields":       string(fields),
		"group_by":     string(groupBy),
	}

	response, err := tool.NewHttp(AccountReportUrl, 5*time.Second).Get(nil, p)

	if err != nil {
		hlog.Error("GDT查询报表数据请求错误1 error: " + err.Error() + " response: " + string(response))
		return err
	}

	res := &AccountReportDataBase{}
	err = json.Unmarshal(response, res)

	if err != nil {
		hlog.Error("GDT查询报表数据解析错误2 error: " + err.Error() + " response: " + string(response))
		return err
	}

	if res.Code != 0 {
		// warningService.NewMediaData().SetPeriod(1, 600).Add("广点通时段数据", string(response))
		hlog.Error("GDT查询报表数据错误3 response: " + string(response))
		return err
	}

	result := &AccountReportDataResponse{}
	err = json.Unmarshal(response, result)

	if err != nil {
		hlog.Error("GDT查询报表数据解析错误4 error: " + err.Error() + " response: " + string(response))
		return err
	}

	*list = result.Data

	return nil
}
