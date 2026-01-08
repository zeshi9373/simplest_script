package basic

type CompanyAddRequest struct {
	CompanyName  string   `json:"company_name" vd:"$!=''"`
	Manager      string   `json:"manager" vd:"$!=''"`
	Username     string   `json:"username" vd:"$!=''"`
	Password     string   `json:"password" vd:"$!=''"`
	Remark       string   `json:"remark"`
	BusinessList []string `json:"business_list"`
}

type CompanyUpdateRequest struct {
	Id           int      `json:"id"`
	CompanyName  string   `json:"company_name" vd:"$!=''"`
	Manager      string   `json:"manager" vd:"$!=''"`
	Remark       string   `json:"remark"`
	BusinessList []string `json:"business_list"`
}

type CompanyPasswordUpdate struct {
	Id              int    `json:"id"`
	NewPassword     string `json:"new_password" vd:"$!=''"`
	ComfirmPassword string `json:"confirm_password" vd:"$!=''"`
}

type CompanyListRequest struct {
	Page        int    `query:"page"`
	PageSize    int    `query:"page_size"`
	CompanyName string `query:"company_name"`
}

type CompanyListResponse struct {
	List     []CompanyListItem `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

type CompanyListItem struct {
	Id           int      `json:"id"`
	CompanyName  string   `json:"company_name"`
	Manager      string   `json:"manager"`
	Remark       string   `json:"remark"`
	Username     string   `json:"username"`
	BusinessList []string `json:"business_list"` // 新增字段
}
