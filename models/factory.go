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
	MessageSession{},
	Message{},
	TranslationSession{},
	PrivateCallSession{},

	//Finance-related
	CoinBundle{},
	Transaction{},
}

func init() {
	log.Println("Initializing Models Factory")

}

//GetModelList will get all models
func GetModelList() []interface{} {
	return modelList
}
