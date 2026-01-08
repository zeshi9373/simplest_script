package export

import (
	"encoding/json"
	"fmt"
	"net/url"
	"simplest_script/core/conf"
	"simplest_script/core/tool"
	"simplest_script/crontab"
	"simplest_script/internal/model/console"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/xuri/excelize/v2"
)

type ApiResponse struct {
	RequestId string       `json:"request_id"`
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	Data      ListResponse `json:"data"`
}

type HeaderItem struct {
	Title string `json:"title"`
	Key   string `json:"key"`
}

type PageInfo struct {
	TotalCount int `json:"total_count"`
	PageSize   int `json:"page_size"`
	Page       int `json:"page"`
	TotalPage  int `json:"total_page"`
}

type ListResponse struct {
	PageInfo PageInfo         `json:"page_info"`
	List     []map[string]any `json:"list"`
}

type Export struct {
}

func (l *Export) Handler(params string) *crontab.Result {
	list := make([]console.ExportLog, 0)
	console.NewExportLogModel().Where("status = 1 and create_time >= ?", time.Now().Add(-24*time.Hour).Format("2006-01-02 15:04:05")).Limit(5).Find(&list)

	for _, log := range list {
		l.Export(log)
	}

	return &crontab.Result{
		Status: 0,
		Data:   nil,
	}
}

func (l *Export) Export(data console.ExportLog) {
	console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
		"status":      2,
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	})

	header := map[string]string{
		"Authorization": "Bearer " + data.Token,
	}

	// 解析URL
	parsedURL, err := url.Parse(data.Query)

	if err != nil {
		console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
			"status":      4,
			"error_msg":   "请求链接解析错误",
			"update_time": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sw, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
			"status":      4,
			"error_msg":   "生成导出文件出错",
			"update_time": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	// 1. 流式写入表头
	headerTitle := []interface{}{}
	dataHeader := make([]HeaderItem, 0)
	err = json.Unmarshal([]byte(data.Header), &dataHeader)
	dataKey := make([]string, 0)
	enums := map[string]map[string]string{}

	for _, v := range dataHeader {
		headerTitle = append(headerTitle, v.Title)
		dataKey = append(dataKey, v.Key)
	}

	json.Unmarshal([]byte(data.Enums), &enums)

	cellAddr, _ := excelize.CoordinatesToCellName(1, 1) // 计算 A1 的坐标
	if err := sw.SetRow(cellAddr, headerTitle); err != nil {
		console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
			"status":      4,
			"error_msg":   "写入表头失败",
			"update_time": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	rows := 2
	page := 1
	pageSize := 500

	for {
		// 获取查询参数
		query := parsedURL.Query()

		// 替换参数值
		query.Set("page", strconv.Itoa(page))          // 将page改为2
		query.Set("page_size", strconv.Itoa(pageSize)) // 将page_size改为20

		// 重新设置查询参数
		parsedURL.RawQuery = query.Encode()
		requestUrl := parsedURL.String()

		if strings.Contains(requestUrl, "%") {
			requestUrl, _ = url.QueryUnescape(requestUrl)
		}

		response, err := tool.NewHttp(requestUrl, 5*time.Second).Get(header, nil)

		if err != nil {
			hlog.Error("查询报表数据请求错误 error: " + err.Error() + " response: " + string(response))
			console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
				"status":      4,
				"update_time": time.Now().Format("2006-01-02 15:04:05"),
			})
			return
		}

		res := &ApiResponse{}
		err = json.Unmarshal(response, res)

		if err != nil {
			hlog.Error("查询报表数据解析错误 error: " + err.Error() + " response: " + string(response))
			console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
				"status":      4,
				"error_msg":   "查询报表数据解析错误 error: " + err.Error() + " response: " + string(response),
				"update_time": time.Now().Format("2006-01-02 15:04:05"),
			})
		}

		if len(res.Data.List) > 0 {
			for _, item := range res.Data.List {
				row := make([]interface{}, 0, len(dataKey))
				for _, k := range dataKey {
					value := item[k]

					if _, ok := enums[k]; ok {
						value = enums[k][fmt.Sprintf("%v", value)]
					}

					row = append(row, value)
				}

				cellAddr, _ := excelize.CoordinatesToCellName(1, rows) // 计算第 i 行第一列的坐标
				if err := sw.SetRow(cellAddr, row); err != nil {
					console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
						"status":      4,
						"error_msg":   "写入报表数据解析错误 error: " + err.Error(),
						"update_time": time.Now().Format("2006-01-02 15:04:05"),
					})
					return
				}

				rows++
			}
		}

		if res.Data.PageInfo.TotalPage <= page {
			break
		}

		page++
	}

	// 保存文件
	filePath := conf.Conf.ExportPath + "/" + data.FileName + time.Now().Format("20060102030405") + strconv.Itoa(tool.Random(10000, 99999)) + ".xlsx"
	if err := f.SaveAs(filePath); err != nil {
		console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
			"status":      4,
			"error_msg":   "保存文件失败",
			"update_time": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	console.NewExportLogModel().Where("id = ?", data.Id).Updates(map[string]any{
		"status":      3,
		"file_path":   filePath,
		"finish_time": time.Now().Format("2006-01-02 15:04:05"),
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	})

}
