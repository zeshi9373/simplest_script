package data

type AccountReportDataRequest struct {
	Dimensions   []string      `json:"dimensions"`
	AdvertiserID string        `json:"advertiser_id"`
	Metrics      []string      `json:"metrics"`
	Filters      []interface{} `json:"filters"`
	Page         string        `json:"page"`
	PageSize     string        `json:"page_size"`
	StartTime    string        `json:"start_time"`
	EndTime      string        `json:"end_time"`
	OrderBy      []interface{} `json:"order_by"`
}

type AccountReportDataBase struct {
	Code int `json:"code"`
}

type AccountReportDataResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Data      Data   `json:"data"`
}
type PageInfo struct {
	Page        int `json:"page"`
	PageSize    int `json:"page_size"`
	TotalNumber int `json:"total_number"`
	TotalPage   int `json:"total_page"`
}
type Dimensions struct {
	StatTimeHour string `json:"stat_time_hour"`
}
type Metrics struct {
	ConvertCnt string `json:"convert_cnt"`
	StatCost   string `json:"stat_cost"`
	ShowCnt    string `json:"show_cnt"`
	ClickCnt   string `json:"click_cnt"`
}
type Rows struct {
	Dimensions Dimensions `json:"dimensions"`
	Metrics    Metrics    `json:"metrics"`
}
type TotalMetrics struct {
	ConvertCnt string `json:"convert_cnt"`
	StatCost   string `json:"stat_cost"`
}
type Data struct {
	PageInfo     PageInfo     `json:"page_info"`
	Rows         []Rows       `json:"rows"`
	TotalMetrics TotalMetrics `json:"total_metrics"`
}
