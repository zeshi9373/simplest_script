package basic

type DepartmentAddRequest struct {
	Name         string `json:"name" vd:"$!=''"`
	ParentId     int    `json:"parent_id"`
	LeaderUserId int    `json:"leader_user_id" vd:"$!=0"`
	Remark       string `json:"remark"`
}

type DepartmentUpdateRequest struct {
	Id           int    `json:"id" vd:"$!=0"`
	Name         string `json:"name"`
	ParentId     int    `json:"parent_id"`
	LeaderUserId int    `json:"leader_user_id" vd:"$!=0"`
	Remark       string `json:"remark"`
}

type DepartmentDeleteRequest struct {
	Id int `json:"id" vd:"$!=0"`
}

type DepartmentListRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Name     string `query:"name"`
}

type DepartmentListItem struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ParentId     int    `json:"parent_id"`
	ParentName   string `json:"parent_name"`
	LeaderUserId int    `json:"leader_user_id"`
	LeaderName   string `json:"leader_name"`
	CompanyId    int    `json:"company_id"`
	Remark       string `json:"remark"`
	State        int    `json:"state"`
}

type DepartmentListResponse struct {
	List     []DepartmentListItem `json:"list"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	Total    int64                `json:"total"`
}

type DepartmentTreeItem struct {
	Id           int                  `json:"id"`
	Name         string               `json:"name"`
	LeaderUserId int                  `json:"leader_user_id"`
	CompanyId    int                  `json:"company_id"`
	Children     []DepartmentTreeItem `json:"children"`
}

type DepartmentTreeListRequest struct {
}

type DepartmentTreeResponse struct {
	List []DepartmentTreeItem `json:"list"`
}
