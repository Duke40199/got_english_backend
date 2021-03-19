package daos

import "log"

var (
	accountDAO    = AccountDAO{}
	coinBundleDAO = CoinBundleDAO{}
)

func init() {
	log.Println("Initializing DAO Factory")

}

func GetAccountDAO() AccountDAO {
	return accountDAO
}
func GetCoinBundleDAO() CoinBundleDAO {
	return coinBundleDAO
}
