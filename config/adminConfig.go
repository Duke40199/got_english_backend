package config

//Config BE model
//RoleNameConfig model
type AdminPerimssionConfig struct {
	CanManageAdmin        string
	CanManageModerator    string
	CanManageLearner      string
	CanManageExpert       string
	CanManageExchangeRate string
}

var permissionConfig = AdminPerimssionConfig{
	CanManageAdmin:        "can_manage_admin",
	CanManageModerator:    "can_manage_moderator",
	CanManageLearner:      "can_manage_learner",
	CanManageExpert:       "can_manage_expert",
	CanManageExchangeRate: "can_manage_exchange_rate",
}

//GetRoleNameConfig : export roleID config
func GetAdminPermissionConfig() *AdminPerimssionConfig {
	return &permissionConfig
}
