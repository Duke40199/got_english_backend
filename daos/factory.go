package daos

import "log"

var (
	adminDAO           = AdminDAO{}
	accountDAO         = AccountDAO{}
	applicationFormDAO = ApplicationFormDAO{}
	coinBundleDAO      = CoinBundleDAO{}
	expertDAO          = ExpertDAO{}
	invoiceDAO         = InvoiceDAO{}
	learnerDAO         = LearnerDAO{}
	moderatorDAO       = ModeratorDAO{}
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
func GetModeratorDAO() ModeratorDAO {
	return moderatorDAO
}
