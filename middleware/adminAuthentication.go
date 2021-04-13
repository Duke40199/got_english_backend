package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
)

func CheckAdminPermission(permission string, r *http.Request) bool {
	var isAuthenticated = false
	switch permission {
	case config.GetPermissionConfig().CanManageAdmin:
		{
			canManageAdmin, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_admin")))
			if canManageAdmin {
				isAuthenticated = true
			}
			break
		}
	case config.GetPermissionConfig().CanManageModerator:
		{
			canManageModerator, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_moderator")))
			if canManageModerator {
				isAuthenticated = true
			}
			break
		}
	case config.GetPermissionConfig().CanManageExpert:
		{
			canManageExpert, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_expert")))
			if canManageExpert {
				isAuthenticated = true
			}
			break
		}
	case config.GetPermissionConfig().CanManageLearner:
		{
			canManageLearner, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_learner")))
			if canManageLearner {
				isAuthenticated = true
			}
			break
		}
	}
	return isAuthenticated
}

func GetPermissionByRoleName(roleName string) string {
	var permission = ""
	switch roleName {
	case roleNameConfig.Admin:
		{
			permission = config.GetPermissionConfig().CanManageAdmin
			break
		}

	case config.GetRoleNameConfig().Moderator:
		{
			permission = config.GetPermissionConfig().CanManageModerator
			break
		}

	case config.GetRoleNameConfig().Expert:
		{
			permission = config.GetPermissionConfig().CanManageExpert
			break
		}

	case config.GetRoleNameConfig().Learner:
		{
			permission = config.GetPermissionConfig().CanManageLearner
			break
		}
	}
	return permission
}
