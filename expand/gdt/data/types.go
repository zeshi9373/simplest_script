package data

type AccountReportDataRequest struct {
	AccountID string    `json:"account_id"`
	Level     string    `json:"level"`
	PageSize  int       `json:"page_size"`
	DateRange DateRange `json:"date_range"`
	Page      int       `json:"page"`
	GroupBy   []string  `json:"group_by"`
	Fields    []string  `json:"fields"`
}
type DateRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type AccountReportDataBase struct {
	Code int `json:"code"`
}

type AccountReportDataResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	MessageCn string `json:"message_cn"`
	Data      Data   `json:"data"`
}
type List struct {
	Cost       int `json:"cost"`
	Hour       int `json:"hour"`
	ConvertCnt int `json:"conversions_count"`
	ViewCount  int `json:"view_count"`
	ClickCount int `json:"valid_click_count"`
}
type PageInfo struct {
	Page        int `json:"page"`
	PageSize    int `json:"page_size"`
	TotalNumber int `json:"total_number"`
	TotalPage   int `json:"total_page"`
}
type Data struct {
	List     []List   `json:"list"`
	PageInfo PageInfo `json:"page_info"`
}
