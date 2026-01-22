package permission

type MenuListAllItem struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	URL   string         `json:"url"`
	Icon  string         `json:"icon"`
	Items []MenuListItem `json:"children"`
}

type MenuListItem struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Icon  string         `json:"icon"`
	URL   string         `json:"url"`
	Items []MenuListItem `json:"children"`
}

type RolePermissionSaveRequest struct {
	RoleId           int            `json:"role_id" vd:"$>0"`
	PermissionIds    map[string]int `json:"permission_id" vd:"$!=''"`
	PathDataTemplate map[string]int `json:"path_data_template" vd:"$!=''"`
}

type RolePermissionDetailRequest struct {
	RoleId int `query:"role_id" vd:"$>0"`
}

type RolePermission struct {
	MenuList   []RoleMenuListItem `json:"menu_list"`
	ButtonList []string           `json:"button_list"`
}

type RoleMenuListItem struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	URL   string         `json:"url"`
	Icon  string         `json:"icon"`
	Items []RoleMenuItem `json:"items"`
}

type RoleMenuItem struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Icon  string         `json:"icon"`
	URL   string         `json:"url"`
	Items []RoleMenuItem `json:"items"`
}
