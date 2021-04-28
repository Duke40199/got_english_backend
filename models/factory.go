package models

import (
	"log"
)

var modelList = []interface{}{
	//Account-related
	Account{},
	Learner{},
	Expert{},
	Admin{},
	Moderator{},
	ApplicationForm{},
	//Finance-related
	Pricing{},
	CoinBundle{},
	Invoice{},
	ExchangeRate{},
	RatingAlgorithm{},
	//Service-related
	MessagingSession{},
	LiveCallSession{},
	Rating{},
}

func init() {
	log.Println("Initializing Models Factory")

}

//GetModelList will get all models
func GetModelList() []interface{} {
	return modelList
}
