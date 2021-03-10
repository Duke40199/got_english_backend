package daos

import "log"

var (
	userDAO = AccountDAO{}
)

func init() {
	log.Println("Initializing DAO Factory")

}

func GetAccountDAO() AccountDAO {
	return userDAO
}
