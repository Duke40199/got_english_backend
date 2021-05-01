package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
)

func CheckModeratorPermission(permission string, r *http.Request) bool {
	var isAuthenticated = false
	switch permission {
	case config.GetModeratorPermissionConfig().CanManageApplicationForm:
		{
			isAuthenticated, _ = strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_application_form")))
			break
		}
	case config.GetModeratorPermissionConfig().CanManageCoinBundle:
		{
			isAuthenticated, _ = strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_coin_bundle")))
			break
		}
	case config.GetModeratorPermissionConfig().CanManageExchangeRate:
		{
			isAuthenticated, _ = strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_exchange_rate")))
			break
		}
	case config.GetModeratorPermissionConfig().CanManagePricing:
		{
			isAuthenticated, _ = strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_pricing")))
			break
		}
	case config.GetModeratorPermissionConfig().CanManageRatingAlgorithm:
		{
			isAuthenticated, _ = strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_rating_algorithm")))
			break
		}
	}
	return isAuthenticated
}
