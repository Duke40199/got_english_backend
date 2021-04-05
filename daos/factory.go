package daos

import "log"

var (
	adminDAO              = AdminDAO{}
	accountDAO            = AccountDAO{}
	applicationFormDAO    = ApplicationFormDAO{}
	coinBundleDAO         = CoinBundleDAO{}
	expertDAO             = ExpertDAO{}
	invoiceDAO            = InvoiceDAO{}
	learnerDAO            = LearnerDAO{}
	messagingSessionDAO   = MessagingSessionDAO{}
	moderatorDAO          = ModeratorDAO{}
	pricingDAO            = PricingDAO{}
	privateCallSessionDAO = PrivateCallSessionDAO{}
	ratingDAO             = RatingDAO{}
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
func GetExpertDAO() ExpertDAO {
	return expertDAO
}
func GetInvoiceDAO() InvoiceDAO {
	return invoiceDAO
}
func GetLearnerDAO() LearnerDAO {
	return learnerDAO
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
func GetPricingDAO() PricingDAO {
	return pricingDAO
}
func GetPrivateCallSessionDAO() PrivateCallSessionDAO {
	return privateCallSessionDAO
}
func GetTranslationSessionDAO() TranslationSessionDAO {
	return translationSessionDAO
}
