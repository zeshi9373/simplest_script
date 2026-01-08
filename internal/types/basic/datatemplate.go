package basic

type DataTemplateAddRequest struct {
	PermissionId int    `json:"permission_id" vd:"$!=''"`
	Title        string `json:"title" vd:"$!=''"`
	HeaderFields string `json:"header_fields" vd:"$!=''"`
}

type DataTemplateUpdateRequest struct {
	Id           int    `json:"id" vd:"$!=0"`
	Title        string `json:"title" vd:"$!=''"`
	HeaderFields string `json:"header_fields" vd:"$!=''"`
}

type DataTemplateDeleteRequest struct {
	Id int `json:"id" vd:"$!=0"`
}

type DataTemplateListRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type DataTemplateListItem struct {
	Id             int    `json:"id"`
	CompanyId      int    `json:"company_id"`
	PermissionId   int    `json:"permission_id"`
	PermissionName string `json:"permission_name"` // 新增字段
	HeaderFields   string `json:"header_fields"`
	Title          string `json:"title"`
	CreateTime     string `json:"create_time"`
}

type DataTemplateListResponse struct {
	List     []DataTemplateListItem `json:"list"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Total    int64                  `json:"total"`
}

type DataTemplateMapItem struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}
type DataTemplateMapResponse map[int][]DataTemplateMapItem
