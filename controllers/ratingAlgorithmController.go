package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/middleware"
	"github.com/golang/got_english_backend/models"
	"github.com/golang/got_english_backend/utils"
)

func GetRatingAlgorithmHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageRatingAlgorithm, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage rating algorithm", http.StatusUnauthorized)
		return
	}

	ratingAlgorithmDAO := daos.GetRatingAlgorithmDAO()
	result, err := ratingAlgorithmDAO.GetRatingAlgorithm()
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func UpdateRatingAlgorithmHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message         = "OK"
		ratingAlgorithm models.RatingAlgorithm
	)
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageRatingAlgorithm, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage rating algorithm", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&ratingAlgorithm); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	if ratingAlgorithm.AverageAllExpertsRating != 0 {
		http.Error(w, "Cannot update average all experts rating", http.StatusBadRequest)
		return
	}
	ratingAlgorithmDAO := daos.GetRatingAlgorithmDAO()
	result, err := ratingAlgorithmDAO.UpdateRatingAlgorithm(1, ratingAlgorithm)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Update expert weighted rating
	//Get all ratings
	ratingDAO := daos.GetRatingDAO()
	ratingList, err := ratingDAO.GetRatings(0)
	//Get expert count
	expertDAO := daos.GetExpertDAO()
	expertRowCount, _ := expertDAO.GetExpertRowCount()
	//Get updated rating algo numbers
	currentRatingAlgorithm, _ := ratingAlgorithmDAO.GetRatingAlgorithm()
	//Update expert rating weighted rating
	for i := 0; i < int(expertRowCount); i++ {
		expertWeightedRating := utils.CalculateExpertWeightedRating(uint(i+1), ratingList, currentRatingAlgorithm)
		_, err = expertDAO.UpdateWeightedRatingByExpertID(uint(i+1), expertWeightedRating)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	}
	config.ResponseWithSuccess(w, message, result)
}
