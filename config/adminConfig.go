package config

//Config BE model
//RoleNameConfig model
type PerimssionConfig struct {
	CanManageAdmin     string
	CanManageModerator string
	CanManageLearner   string
	CanManageExpert    string
}

var permissionConfig = PerimssionConfig{
	CanManageAdmin:     "can_manage_admin",
	CanManageModerator: "can_manage_moderator",
	CanManageLearner:   "can_manage_learner",
	CanManageExpert:    "can_manage_expert",
}

//GetRoleNameConfig : export roleID config
func GetPermissionConfig() *PerimssionConfig {
	return &permissionConfig
}
