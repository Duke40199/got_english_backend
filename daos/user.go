package daos

import (
	"fmt"

	AccountModel "github.com/golang/GotEnglishBackend/Application/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (u *AccountDAO) CreateUser(db *gorm.DB, user AccountModel.Account) (*AccountModel.Account, error) {
	err := db.Debug().Create(&user).Error
	if err != nil {
		return &AccountModel.Account{}, err
	}
	return nil, nil
}

func (u *AccountDAO) FindUserByUsername(db *gorm.DB, user AccountModel.Account) (*AccountModel.Account, error) {
	err := db.Where("username LIKE ?", user.Username).First(&user).Error
	if err == nil {
		return &user, err
	}
	return nil, nil
}

func (u *AccountDAO) FindUserByUsernameAndPassword(db *gorm.DB, user AccountModel.Account) (*AccountModel.Account, error) {
	var result = AccountModel.Account{}
	err := db.Debug().First(&result, "username=?", user.Username).Error
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
		if err != nil {
			return &AccountModel.Account{}, err
		}
		return &result, err
	}
	return &AccountModel.Account{}, nil
}

func (u *AccountDAO) FindUserByEmailAndPassword(db *gorm.DB, user AccountModel.Account) (*AccountModel.Account, error) {
	var result = AccountModel.Account{}
	err := db.Debug().First(&result, "email=?", user.Email).Error
	if err == nil {
		fmt.Println("====== ERROR NOT FOUND")
		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
		if err != nil {
			fmt.Println("====== ERROR FOUND")
			return &AccountModel.Account{}, err
		}
		return &result, err
	}
	return &AccountModel.Account{}, nil
}
