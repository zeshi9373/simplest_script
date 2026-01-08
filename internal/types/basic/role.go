package basic

type RoleAddRequest struct {
	Name   string `json:"name" vd:"$!=''"`
	Remark string `json:"remark"`
}

type RoleUpdateRequest struct {
	Id     int    `json:"id" vd:"$>0"`
	Name   string `json:"name" vd:"$!=''"`
	Remark string `json:"remark"`
}

type RoleListRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type RoleListResponse struct {
	List     []RoleListItem `json:"list"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Total    int64          `json:"total"`
}

type RoleListItem struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Remark       string `json:"remark"`
	PermissionId string `json:"permission_id"`
	CreateTime   string `json:"create_time"`
}
