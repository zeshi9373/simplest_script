package basic

type AdminUserAddRequest struct {
	Username string `json:"username" vd:"$!=''"` // 用户名必填
	Password string `json:"password" vd:"$!=''"` // 密码必填
	Nickname string `json:"nickname" vd:"$!=''"` // 姓名必填
	DepartId int    `json:"depart_id"`           // 必填，手动判断 !=0
	Status   int    `json:"status"`              // 状态，可选
	RoleIds  string `json:"role_ids"`            // 角色ID，可选
	Remark   string `json:"remark"`              // 备注，可选
}

type AdminUserUpdateRequest struct {
	Id       int    `json:"id"`        // 必填
	Nickname string `json:"nickname"`  // 可选
	DepartId int    `json:"depart_id"` // 可选
	Status   int    `json:"status"`    // 可选
	RoleIds  string `json:"role_ids"`  // 可选
	Remark   string `json:"remark"`    // 可选
}

type AdminUserPasswordUpdate struct {
	Id              int    `json:"id"` // 必填
	NewPassword     string `json:"new_password" vd:"$!=''"`
	ComfirmPassword string `json:"confirm_password" vd:"$!=''"`
}

type AdminUserDeleteRequest struct {
	Id int `json:"id"` // 必填
}

type AdminUserListRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	DepartId int    `query:"depart_id"`
	RoleId   int    `query:"role_id"`
	Status   int    `query:"status"`
	Username string `query:"username"`
	Nickname string `query:"nickname"`
}

type AdminUserListItem struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	DepartId   int    `json:"depart_id"`
	DepartName string `json:"depart_name"`
	Status     int    `json:"status"`
	RoleIds    string `json:"role_ids"`
	RoleName   string `json:"role_name"`
	Remark     string `json:"remark"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type AdminUserListResponse struct {
	List     []AdminUserListItem `json:"list"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
	Total    int64               `json:"total"`
}
