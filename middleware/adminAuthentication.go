package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/got_english_backend/config"
)

func CheckAdminPermission(permission string, r *http.Request) bool {
	var isAuthenticated = false
	switch permission {
	case config.GetAdminPermissionConfig().CanManageAdmin:
		{
			canManageAdmin, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_admin")))
			if canManageAdmin {
				isAuthenticated = true
			}
			break
		}
	case config.GetAdminPermissionConfig().CanManageModerator:
		{
			canManageModerator, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_moderator")))
			if canManageModerator {
				isAuthenticated = true
			}
			break
		}
	case config.GetAdminPermissionConfig().CanManageExpert:
		{
			canManageExpert, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_expert")))
			if canManageExpert {
				isAuthenticated = true
			}
			break
		}
	case config.GetAdminPermissionConfig().CanManageLearner:
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

func GetAdminPermissionByRoleName(roleName string) string {
	var permission = ""
	switch strings.Title(roleName) {
	case roleNameConfig.Admin:
		{
			permission = config.GetAdminPermissionConfig().CanManageAdmin
			break
		}

	case config.GetRoleNameConfig().Moderator:
		{
			permission = config.GetAdminPermissionConfig().CanManageModerator
			break
		}

	case config.GetRoleNameConfig().Expert:
		{
			permission = config.GetAdminPermissionConfig().CanManageExpert
			break
		}

	case config.GetRoleNameConfig().Learner:
		{
			permission = config.GetAdminPermissionConfig().CanManageLearner
			break
		}
	}
	return permission
}
