package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/utils"
)

func GetExpertEarningsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message  = "OK"
		duration = "weekly"
	)
	var err error
	//Get duration
	if len(r.URL.Query()["duration"]) > 0 {
		duration = r.URL.Query()["duration"][0]
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	}
	startDate, endDate := utils.GetTimesByPeriod(duration)
	expertID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("expert_id")), 10, 0)
	earningDAO := daos.GetEarningDAO()
	result, err := earningDAO.GetEarningByExpertID(uint(expertID), startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
