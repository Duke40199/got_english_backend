package daos

import "log"

var (
	adminDAO              = AdminDAO{}
	accountDAO            = AccountDAO{}
	applicationFormDAO    = ApplicationFormDAO{}
	coinBundleDAO         = CoinBundleDAO{}
	earningDAO            = EarningDAO{}
	exchangeRateDAO       = ExchangeRateDAO{}
	expertDAO             = ExpertDAO{}
	invoiceDAO            = InvoiceDAO{}
	learnerDAO            = LearnerDAO{}
	liveCallSessionDAO    = LiveCallSessionDAO{}
	messagingSessionDAO   = MessagingSessionDAO{}
	moderatorDAO          = ModeratorDAO{}
	pricingDAO            = PricingDAO{}
	ratingDAO             = RatingDAO{}
	ratingAlgorithmDAO    = RatingAlgorithmDAO{}
	translationSessionDAO = TranslationSessionDAO{}
)

func init() {
	log.Println("Initializing DAO Factory")

}
func GetAdminDAO() AdminDAO {
	return adminDAO
}
func GetAccountDAO() AccountDAO {
	return accountDAO
}
func GetApplicationFormDAO() ApplicationFormDAO {
	return applicationFormDAO
}
func GetCoinBundleDAO() CoinBundleDAO {
	return coinBundleDAO
}
func GetEarningDAO() EarningDAO {
	return earningDAO
}
func GetExchangeRateDAO() ExchangeRateDAO {
	return exchangeRateDAO
}
func GetExpertDAO() ExpertDAO {
	return expertDAO
}
func GetInvoiceDAO() InvoiceDAO {
	return invoiceDAO
}
func GetLearnerDAO() LearnerDAO {
	return learnerDAO
}
func GetLiveCallSessionDAO() LiveCallSessionDAO {
	return liveCallSessionDAO
}
func GetMessagingSessionDAO() MessagingSessionDAO {
	return messagingSessionDAO
}
func GetModeratorDAO() ModeratorDAO {
	return moderatorDAO
}
func GetRatingDAO() RatingDAO {
	return ratingDAO
}
func GetRatingAlgorithmDAO() RatingAlgorithmDAO {
	return ratingAlgorithmDAO
}

func GetPricingDAO() PricingDAO {
	return pricingDAO
}

func GetTranslationSessionDAO() TranslationSessionDAO {
	return translationSessionDAO
}
