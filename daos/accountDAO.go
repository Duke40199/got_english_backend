package daos

import (
	"errors"
	"strings"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type AccountDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (u *AccountDAO) CreateAccount(account models.Account) (*models.Account, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Omit("birthday").Create(&account).Error
	return &account, err

}

func (u *AccountDAO) FindUserByUsername(account models.Account) (*models.AccountFullInfo, error) {
	accountResult := models.AccountFullInfo{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Select("accounts.*, experts.can_chat, experts.can_join_translation_session, experts.can_join_private_call_session, admins.can_manage_expert,admins.can_manage_learner,admins.can_manage_admin, admins.can_manage_moderator, moderators.can_manage_coin_bundle,moderators.can_manage_pricing,moderators.can_manage_application_form").
		Where("accounts.username = ?", account.Username).
		Joins("left join experts on experts.account_id = accounts.id").
		Joins("left join learners on learners.account_id = accounts.id").
		Joins("left join moderators on moderators.account_id = accounts.id").
		Joins("left join admins on admins.account_id = accounts.id").
		First(&accountResult).Error

	//Only get date from birthdays
	accountResult.Birthday = strings.Split(accountResult.Birthday, "T")[0]
	GetRoleIDByAccountID(accountResult.RoleName, &accountResult)
	return &accountResult, err
}

func (u *AccountDAO) FindAccountByUsernameAndPassword(account models.Account) (*models.Account, error) {
	var result = models.Account{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().First(&result, "username=?", account.Username).Error
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(*result.Password), []byte(*account.Password))
		if err != nil {
			return &models.Account{}, err
		}
	}
	return &result, nil
}

func (u *AccountDAO) FindAccountByEmailAndPassword(account models.Account) (*models.Account, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Account{}
	err = db.Debug().First(&result, "email=?", account.Email).Error
	if err == nil {
		//If someone logs into Google but haven't updated their password
		if result.Password == nil || *result.Password == "" {
			return &models.Account{}, errors.New("are you logging from google? if so, have you updated your password?")
		}
		err := bcrypt.CompareHashAndPassword([]byte(*result.Password), []byte(*account.Password))
		if err != nil {
			return &models.Account{}, err
		}
	}
	return &result, nil
}

func (u *AccountDAO) FindAccountByEmail(account models.Account) (*models.Account, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Account{}
	err = db.Debug().First(&result, "email=?", account.Email).Error
	return &result, err
}

func (u *AccountDAO) GetAccounts(account models.Account) (*[]models.AccountFullInfo, error) {
	accounts := []models.AccountFullInfo{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Select("accounts.*, experts.can_chat, experts.can_join_translation_session, experts.can_join_private_call_session, admins.can_manage_expert,admins.can_manage_learner,admins.can_manage_admin, admins.can_manage_moderator, moderators.can_manage_coin_bundle,moderators.can_manage_pricing,moderators.can_manage_application_form").
		Where("accounts.role_name LIKE ? AND accounts.username LIKE ?", account.RoleName+"%", *account.Username+"%").
		Joins("left join experts on experts.account_id = accounts.id").
		Joins("left join learners on learners.account_id = accounts.id").
		Joins("left join moderators on moderators.account_id = accounts.id").
		Joins("left join admins on admins.account_id = accounts.id").
		Scan(&accounts).Error
	//Only get date from birthdays
	for i := 0; i < len(accounts); i++ {
		accounts[i].Birthday = strings.Split(accounts[i].Birthday, "T")[0]
	}
	return &accounts, err
}
func (u *AccountDAO) UpdateAccountByID(accountID uuid.UUID, updateInfo map[string]interface{}) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Account{}).Where("id = ?", accountID).
		Updates(updateInfo)
	return result.RowsAffected, result.Error
}

func GetRoleIDByAccountID(roleName string, accountFullInfo *models.AccountFullInfo) {
	switch roleName {
	case config.GetRoleNameConfig().Learner:
		{
			learnerDAO := GetLearnerDAO()
			learner, err := learnerDAO.GetLearnerInfoByAccountID(accountFullInfo.ID)
			if err == nil {
				accountFullInfo.LearnerID = learner.ID
			}
			break
		}
	case config.GetRoleNameConfig().Expert:
		{
			expertDAO := GetExpertDAO()
			expert, err := expertDAO.GetExpertByAccountID(accountFullInfo.ID)
			if err == nil {
				accountFullInfo.ExpertID = expert.ID
			}
			break
		}
	case config.GetRoleNameConfig().Admin:
		{
			adminDAO := GetAdminDAO()
			admin, err := adminDAO.GetAdminByAccountID(accountFullInfo.ID)
			if err == nil {
				accountFullInfo.LearnerID = admin.ID
			}
			break
		}
	case config.GetRoleNameConfig().Moderator:
		{
			moderatorDAO := GetModeratorDAO()
			moderator, err := moderatorDAO.GetModeratorByAccountID(accountFullInfo.ID)
			if err == nil {
				accountFullInfo.ModeratorID = moderator.ID
			}
			break
		}
	}
}
