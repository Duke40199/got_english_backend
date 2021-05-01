package config

//Config BE model
//RoleNameConfig model
type ModeratorPermissionConfig struct {
	CanManageCoinBundle      string
	CanManagePricing         string
	CanManageApplicationForm string
	CanManageExchangeRate    string
	CanManageRatingAlgorithm string
}

var moderatorPermissionConfig = ModeratorPermissionConfig{
	CanManageCoinBundle:      "can_manage_coin_bundle",
	CanManagePricing:         "can_manage_pricing",
	CanManageApplicationForm: "can_manage_application_form",
	CanManageExchangeRate:    "can_manage_exchange_rate",
	CanManageRatingAlgorithm: "can_manage_rating_algorithm",
}

//GetRoleNameConfig : export roleID config
func GetModeratorPermissionConfig() *ModeratorPermissionConfig {
	return &moderatorPermissionConfig
}
