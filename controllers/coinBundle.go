package controllers

import (
	"encoding/json"
	"net/http"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
)

func CreateCoinBundleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	var coinBundle = models.CoinBundle{}

	if err := json.NewDecoder(r.Body).Decode(&coinBundle); err != nil {
		errMsg := "Malformed data"
		responseConfig.ResponseWithError(w, errMsg, err)
	}
	coinBundleDAO := daos.GetCoinBundleDAO()
	result, err := coinBundleDAO.CreateCoinBundle(models.CoinBundle{
		Title:       coinBundle.Title,
		Description: coinBundle.Description,
		Quantity:    uint(coinBundle.Quantity),
		Price:       uint(coinBundle.Price),
		PriceUnit:   coinBundle.PriceUnit,
	},
	)
	if err != nil {
		panic(err)
	} else {
		responseConfig.ResponseWithSuccess(w, message, result)
	}
}
