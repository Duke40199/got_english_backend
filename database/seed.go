package database

import (
	"fmt"

	"github.com/golang/GotEnglishBackend/Application/config"
	"github.com/golang/GotEnglishBackend/Application/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var roleNameConfig = config.GetRoleNameConfig()

//SeedDB function will trigger all seed functions below
func SeedDB(db *gorm.DB) {

	accounts := []models.Account{
		{
			ID:       uuid.New(),
			Username: "SpacePotato",
			Password: "password",
			Email:    "anhntse130266@fpt.edu.vn",
			RoleName: roleNameConfig.Admin,
		},
		{
			ID:       uuid.New(),
			Username: "DucPhi",
			Password: "password",
			Email:    "duc@phi.com",
			RoleName: roleNameConfig.Expert,
		},
		{
			ID:       uuid.New(),
			Username: "TuanNguyen",
			Password: "password",
			Email:    "tuan@nguyen.com",
			RoleName: roleNameConfig.Moderator,
		},
		{
			ID:       uuid.New(),
			Username: "LocTr",
			Password: "password",
			Email:    "loc@tr.com",
			RoleName: roleNameConfig.Learner,
		},
	}
	SeedAccounts(db, &accounts)
	SeedRolesForAccounts(db, &accounts)
}

//SeedAccounts will seed users to the DB
func SeedAccounts(db *gorm.DB, accounts *[]models.Account) {
	db.Create(&accounts)
	fmt.Println("======= User seeded.")
}

//SeedRolesForAccounts will seed users to the DB
func SeedRolesForAccounts(db *gorm.DB, accounts *[]models.Account) {
	for i := 0; i < len(*accounts); i++ {
		switch (*accounts)[i].RoleName {
		case roleNameConfig.Admin:
			{
				db.Create(&models.Admin{
					AccountsID:       (*accounts)[i].ID,
					CanManageExpert:  true,
					CanManageLearner: true,
					CanManageAdmin:   true,
				})
				break
			}
		case roleNameConfig.Learner:
			{
				db.Create(&models.Learner{
					AccountsID: (*accounts)[i].ID,
				})
				break
			}
		case roleNameConfig.Expert:
			{
				db.Create(&models.Expert{
					AccountsID:                (*accounts)[i].ID,
					CanChat:                   true,
					CanJoinTranslationSession: true,
					CanJoinPrivateCallSession: true,
				})
				break
			}
		case roleNameConfig.Moderator:
			{
				db.Create(&models.Moderator{
					AccountsID:               (*accounts)[i].ID,
					CanManageCoinBundle:      true,
					CanManagePricing:         true,
					CanManageApplicationForm: true,
				})
				break
			}
		}
	}
	fmt.Println("======= User seeded.")
}
