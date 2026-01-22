package user

type OperationLoginRequest struct {
	Email    string `json:"email" vd:"$!=''"`
	Password string `json:"password"  vd:"$!=''"`
}

type OperationLoginResponse struct {
	Token            string              `json:"token"`
	UserId           int                 `json:"user_id"`
	Username         string              `json:"username"`
	Nickname         string              `json:"nickname"`
	CompanyId        int                 `json:"company_id"`
	DepartId         int                 `json:"depart_id"`
	RoleIds          string              `json:"role_ids"`
	Remark           string              `json:"remark"`
	IsAdmin          int                 `json:"is_admin"`
	IsManage         int                 `json:"is_manage"`
	BtnPermission    []string            `json:"btn_permission"`
	PathUrl          string              `json:"path_url"`
	PathDataTemplate map[int]string      `json:"path_data_template"`
	Enums            map[string][]string `json:"enums"`
}
