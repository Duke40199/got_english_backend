package daos

import (
	"strings"
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"

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

type AccountFullInfo struct {
	gorm.Model `json:"-"`
	//Login
	ID       uuid.UUID `gorm:"size:255;column:id;not null;unique; primaryKey;" json:"id"`
	Username string    `gorm:"size:255;not null;unique" json:"username"`
	Fullname string    `gorm:"size:255;not null;unique" json:"fullname"`
	Email    string    `gorm:"size:100;not null;unique" json:"email"`
	Password string    `gorm:"size:100;not null;" json:"-"`
	RoleName string    `gorm:"size:100;not null;" json:"role_name"`
	//Info
	AvatarURL   string     `gorm:"size:255" json:"avatar_url"`
	Address     string     `gorm:"size:255;" json:"address"`
	PhoneNumber string     `gorm:"column:phone_number;autoCreateTime" json:"phone_number"`
	Birthday    string     `gorm:"column:birthday;type:date" json:"birthday" sql:"date"`
	IsSuspended bool       `gorm:"column:isSuspended" json:"is_suspended"`
	SuspendedAt *time.Time `gorm:"column:SuspendedAt" json:"suspended_at"`
	//default timestamps
	CreatedAt time.Time  `gorm:"column:CreatedAt;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:UpdatedAt;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:DeletedAt" json:"deleted_at"`
	//expert perms
	CanChat                   bool `gorm:"column:can_chat" json:"can_chat"`
	CanJoinTranslationSession bool `gorm:"column:can_join_translation_session" json:"can_join_translation_session"`
	CanJoinPrivateCallSession bool `gorm:"column:can_private_call_session" json:"can_private_call_session"`
	//admin perms
	CanManageExpert  bool `gorm:"column:can_manage_expert" json:"can_manage_expert"`
	CanManageLearner bool `gorm:"column:can_manage_learner" json:"can_manage_learner"`
	CanManageAdmin   bool `gorm:"column:can_manage_admin" json:"can_manage_admin"`
	//moderator perms
	CanManageCoinBundle      bool `gorm:"column:can_manage_coin_bundle" json:"can_manage_coin_bundle"`
	CanManagePricing         bool `gorm:"column:can_manage_pricing" json:"can_manage_pricing"`
	CanManageApplicationForm bool `gorm:"column:can_manage_application_form" json:"can_manage_application_form"`
}

func (u *AccountDAO) CreateAccount(account models.Account) (*models.Account, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Omit("birthday").Create(&account).Error
	return &account, err

}

func (u *AccountDAO) FindUserByUsername(account models.Account) (*AccountFullInfo, error) {
	accountResult := AccountFullInfo{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Select("accounts.*, experts.can_chat, experts.can_join_translation_session, experts.can_join_private_call_session, admins.can_manage_expert,admins.can_manage_learner,admins.can_manage_admin, moderators.can_manage_coin_bundle,moderators.can_manage_pricing,moderators.can_manage_application_form").
		Where("accounts.username = ?", account.Username).
		Joins("left join experts on experts.account_id = accounts.id").
		Joins("left join learners on learners.account_id = learners.id").
		Joins("left join moderators on moderators.account_id = accounts.id").
		Joins("left join admins on admins.account_id = accounts.id").
		First(&accountResult).Error
	//Only get date from birthdays
	accountResult.Birthday = strings.Split(accountResult.Birthday, "T")[0]

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
		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(account.Password))
		if err != nil {
			return &models.Account{}, err
		}
		return &result, err
	}
	return &models.Account{}, nil
}

func (u *AccountDAO) FindAccountByEmailAndPassword(account models.Account) (*models.Account, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Account{}
	err = db.Debug().First(&result, "email=?", account.Email).Error
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(account.Password))
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
}

func (u *AccountDAO) GetAccounts(account models.Account) (*[]AccountFullInfo, error) {
	accounts := []AccountFullInfo{}
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Account{}).Select("accounts.*, experts.can_chat, experts.can_join_translation_session, experts.can_join_private_call_session, admins.can_manage_expert,admins.can_manage_learner,admins.can_manage_admin, moderators.can_manage_coin_bundle,moderators.can_manage_pricing,moderators.can_manage_application_form").
		Where("accounts.role_name LIKE ? AND accounts.username LIKE ?", account.RoleName+"%", account.Username+"%").
		Joins("left join experts on experts.account_id = accounts.id").
		Joins("left join learners on learners.account_id = learners.id").
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
