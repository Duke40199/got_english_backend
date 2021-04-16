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

	//Service-related
	MessagingSession{},
	TranslationSession{},
	LiveCallSession{},
	Rating{},

	//Finance-related
	Pricing{},
	CoinBundle{},
	Invoice{},
	ExchangeRate{},
	Earning{},
}

func init() {
	log.Println("Initializing Models Factory")

}

//GetModelList will get all models
func GetModelList() []interface{} {
	return modelList
}
