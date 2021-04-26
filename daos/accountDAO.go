package daos

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/golang/got_english_backend/utils"
	"github.com/google/uuid"
)

type AccountDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (u *AccountDAO) CreateAccount(account models.Account, permissions models.PermissionStruct) (*models.Account, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	//validate if account already exists
	accountAvailable, _ := accountDAO.FindAccountByEmail(account)
	if accountAvailable.Email != nil {
		return &account, errors.New("account unavailable.")
	}
	//Generate username if reqbody doesn't have one
	if account.Username == nil {
		currentTimeMillis := utils.GetCurrentEpochTimeInMiliseconds()
		newUsername := account.RoleName + strconv.FormatInt(currentTimeMillis, 10)
		account.Username = &newUsername
	}
	//if id is not inputted, generate one.
	if account.ID.String() == "00000000-0000-0000-0000-000000000000" {
		account.ID = uuid.New()
	}
	err = db.Debug().Model(&models.Account{}).Create(&account).Error
	if err != nil {
		return nil, err
	}
	switch account.RoleName {
	case config.GetRoleNameConfig().Learner:
		{
			err = db.Debug().Model(&account).Association("Learner").Append(&models.Learner{
				AvailableCoinCount: 0,
			})
			break
		}
	case config.GetRoleNameConfig().Expert:
		{
			err = db.Debug().Model(&account).Association("Expert").Append(&models.Expert{
				CanChat:                   permissions.CanChat,
				CanJoinTranslationSession: permissions.CanJoinTranslationSession,
				CanJoinLiveCallSession:    permissions.CanJoinLiveCallSession,
			})
			break
		}
	case config.GetRoleNameConfig().Moderator:
		{
			err = db.Debug().Model(&account).Association("Moderator").Append(&models.Moderator{
				CanManageCoinBundle:      permissions.CanManageCoinBundle,
				CanManagePricing:         permissions.CanManagePricing,
				CanManageApplicationForm: permissions.CanManageApplicationForm,
			})
			break
		}
	case config.GetRoleNameConfig().Admin:
		{
			err = db.Debug().Model(&account).Association("Admin").Append(&models.Admin{
				CanManageExpert:    permissions.CanManageExpert,
				CanManageLearner:   permissions.CanManageLearner,
				CanManageAdmin:     permissions.CanManageAdmin,
				CanManageModerator: permissions.CanManageAdmin,
			})
			break
		}
	}
	return &account, err
}

func (u *AccountDAO) FindAccountByID(id uuid.UUID) (*models.Account, error) {
	accountResult := models.Account{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Preload("Learner").Preload("Expert").Preload("Moderator").Preload("Admin").
		First(&accountResult, "id=?", fmt.Sprint(id)).Error
	if err != nil {
		return nil, errors.New("account not found.")
	}
	//Only get date from birthdays
	if accountResult.Birthday != nil {
		*accountResult.Birthday = strings.Split(*accountResult.Birthday, "T")[0]
	}
	return &accountResult, err
}

func (u *AccountDAO) FindAccountByUsername(account models.Account) (*models.Account, error) {
	accountResult := models.Account{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Preload("Learner").Preload("Expert").Preload("Moderator").Preload("Admin").
		First(&accountResult, "username=?", account.Username).Error
	//Only get date from birthdays
	if accountResult.Birthday != nil {
		*accountResult.Birthday = strings.Split(*accountResult.Birthday, "T")[0]
	}
	return &accountResult, err
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

func (u *AccountDAO) GetAccounts(account models.Account) (*[]models.Account, error) {
	accounts := []models.Account{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Preload("Learner").Preload("Expert").Preload("Moderator").Preload("Admin").
		Find(&accounts, "accounts.role_name LIKE ? AND accounts.username LIKE ?", "%"+account.RoleName+"%", "%"+*account.Username+"%").Error
	//Only get date from birthdays
	for i := 0; i < len(accounts); i++ {
		if accounts[i].Birthday != nil {
			*accounts[i].Birthday = strings.Split(*accounts[i].Birthday, "T")[0]
		}
		if accounts[i].RoleName == config.GetRoleNameConfig().Expert {
			accounts[i].Expert.AverageRating, err = ratingDAO.GetExpertAverageRating(*accounts[i].Expert)
		}
	}
	return &accounts, err
}
func (u *AccountDAO) GetAccountByAccountID(id uuid.UUID) (*models.Account, error) {
	account := models.Account{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).
		Preload("Learner").
		Preload("Learner.ID").
		Preload("Expert").
		Preload("Expert.ID").
		Preload("Moderator").
		Preload("Admin").
		Find(&account, "id=?", id).Error
	//Only get date from birthdays
	if account.Birthday != nil {
		*account.Birthday = strings.Split(*account.Birthday, "T")[0]
	}
	if account.RoleName == config.GetRoleNameConfig().Expert {
		account.Expert.AverageRating, err = ratingDAO.GetExpertAverageRating(*account.Expert)
	}
	return &account, err
}
func (u *AccountDAO) UpdateAccountByID(accountID uuid.UUID, updateInfo models.Account) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Account{}).Where("id = ?", accountID).
		Updates(&updateInfo)
	return result.RowsAffected, result.Error
}

func (u *AccountDAO) SuspendAccountByID(accountID uuid.UUID) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Account{}).Where("id = ?", accountID).
		Updates(map[string]interface{}{"is_suspended": true, "suspended_at": time.Now()})
	return result.RowsAffected, result.Error
}

func (u *AccountDAO) UnsuspendAccountByID(accountID uuid.UUID) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Account{}).Where("id = ?", accountID).
		Updates(map[string]interface{}{"is_suspended": false, "suspended_at": nil})
	return result.RowsAffected, result.Error
}
