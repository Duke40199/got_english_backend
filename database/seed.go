package database

import (
	"fmt"
	"time"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var roleNameConfig = config.GetRoleNameConfig()

//SeedDB function will trigger all seed functions below
func SeedDB(db *gorm.DB) {

	accounts := []models.Account{
		{
			ID:          uuid.New(),
			Username:    "SpacePotato",
			Password:    "password",
			Email:       "anhntse130266@fpt.edu.vn",
			AvatarURL:   "https://scontent-hkt1-1.xx.fbcdn.net/v/t1.0-9/136993742_3749756998379656_496491351021530940_n.jpg?_nc_cat=107&ccb=1-3&_nc_sid=09cbfe&_nc_ohc=kvndP9qcIusAX-iZQVs&_nc_ht=scontent-hkt1-1.xx&oh=5768c6a74d45c1c23b69b2ef86dfa77c&oe=6072A896",
			Address:     "139 Lac Long Quan, Ward 10,Tan Binh Dist.,Ho Chi Minh City",
			Birthday:    time.Date(1987, 1, 12, 0, 0, 0, 0, time.Now().Location()),
			PhoneNumber: "0123456789",
			RoleName:    roleNameConfig.Admin,
		},
		{
			ID:          uuid.New(),
			Username:    "TuanAnh",
			Password:    "password",
			Email:       "binguyentuananh@gmail.com",
			AvatarURL:   "https://scontent-hkt1-1.xx.fbcdn.net/v/t1.0-9/136993742_3749756998379656_496491351021530940_n.jpg?_nc_cat=107&ccb=1-3&_nc_sid=09cbfe&_nc_ohc=kvndP9qcIusAX-iZQVs&_nc_ht=scontent-hkt1-1.xx&oh=5768c6a74d45c1c23b69b2ef86dfa77c&oe=6072A896",
			Address:     "139 Lac Long Quan, Ward 10,Tan Binh Dist.,Ho Chi Minh City",
			Birthday:    time.Date(1987, 1, 12, 0, 0, 0, 0, time.Now().Location()),
			PhoneNumber: "0123456789",
			RoleName:    roleNameConfig.Moderator,
		},
		{
			ID:          uuid.New(),
			Username:    "DucPhi",
			Password:    "password",
			Email:       "duc@phi.com",
			Address:     "1722  Cody Ridge Road, Enid, Oklahoma",
			Birthday:    time.Date(1999, 12, 12, 0, 0, 0, 0, time.Now().Location()),
			PhoneNumber: "0987654321",
			RoleName:    roleNameConfig.Expert,
		},
		{
			ID:          uuid.New(),
			Username:    "TuanNguyen",
			Password:    "password",
			Email:       "tuan@nguyen.com",
			Birthday:    time.Date(1998, 11, 12, 0, 0, 0, 0, time.Now().Location()),
			Address:     "10/2 Dang Van Ngu St., Ward 10, Phu Nhuan Dist.,Ho Chi Minh City ",
			PhoneNumber: "0777984632",
			RoleName:    roleNameConfig.Admin,
		},
		{
			ID:          uuid.New(),
			Username:    "LocTr",
			Password:    "password",
			Email:       "loc@tr.com",
			Birthday:    time.Date(1999, 12, 12, 0, 0, 0, 0, time.Now().Location()),
			Address:     "4668  Delaware Avenue, San Francisco, California ",
			PhoneNumber: "0334433221",
			RoleName:    roleNameConfig.Learner,
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
					AccountID:        (*accounts)[i].ID,
					CanManageExpert:  true,
					CanManageLearner: true,
					CanManageAdmin:   true,
				})
				break
			}
		case roleNameConfig.Learner:
			{
				db.Create(&models.Learner{
					AccountID: (*accounts)[i].ID,
				})
				break
			}
		case roleNameConfig.Expert:
			{
				db.Create(&models.Expert{
					AccountID:                 (*accounts)[i].ID,
					CanChat:                   true,
					CanJoinTranslationSession: true,
					CanJoinPrivateCallSession: true,
				})
				break
			}
		case roleNameConfig.Moderator:
			{
				db.Create(&models.Moderator{
					AccountID:                (*accounts)[i].ID,
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
