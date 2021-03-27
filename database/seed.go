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
	accounts := []models.Account{}
	usernames := []string{"TuanAnh", "DucPhi", "TuanNguyen", "LocTr"}
	fullnames := []string{"Nguyen Tuan Anh", "Phi Do Hong Duc", "Nguyen Tran Quoc Tuan", "Tran Thien Loc"}
	passwords := []string{"password", "password", "password", "password"}
	emails := []string{"anhntse130266@fpt.edu.vn", "hongduc5412@gmail.com", "tuanntqse62351@fpt.edu.vn", "ttloc1999@gmail.com"}
	avatarUrls := []string{
		"https://scontent-hkt1-1.xx.fbcdn.net/v/t1.0-9/136993742_3749756998379656_496491351021530940_n.jpg?_nc_cat=107&ccb=1-3&_nc_sid=09cbfe&_nc_ohc=kvndP9qcIusAX-iZQVs&_nc_ht=scontent-hkt1-1.xx&oh=5768c6a74d45c1c23b69b2ef86dfa77c&oe=6072A896",
		"",
		"https://lh3.googleusercontent.com/a-/AOh14GjUZVucuFzLCa25tllc6tf2Oh8DZAr32hvbTXl5=s256",
		"https://lh3.googleusercontent.com/a-/AOh14GignIoa0zabTMuuTD2iwR8H_Ph4KZBIReiXdOyF=s256",
	}
	phoneNumbers := []string{"+14155552671", "+44207183875044", "+8477984632", "+843344221"}
	addresses := []string{
		"139 Lac Long Quan, Ward 10,Tan Binh Dist.,Ho Chi Minh City",
		"1722  Cody Ridge Road, Enid, Oklahoma",
		"10/2 Dang Van Ngu St., Ward 10, Phu Nhuan Dist.,Ho Chi Minh City ",
		"4668  Delaware Avenue, San Francisco, California ",
	}
	roleNames := []string{
		roleNameConfig.Moderator,
		roleNameConfig.Expert,
		roleNameConfig.Admin,
		roleNameConfig.Learner,
	}
	birthdays := []time.Time{
		time.Date(1987, 1, 12, 0, 0, 0, 0, time.Now().Location()),
		time.Date(1998, 04, 03, 0, 0, 0, 0, time.Now().Location()),
		time.Date(1997, 11, 14, 0, 0, 0, 0, time.Now().Location()),
		time.Date(1996, 05, 14, 0, 0, 0, 0, time.Now().Location()),
	}
	for i := 0; i < len(usernames); i++ {
		accounts = append(accounts,
			models.Account{
				ID:          uuid.New(),
				Username:    &usernames[i],
				Password:    &passwords[i],
				Fullname:    &fullnames[i],
				Email:       &emails[i],
				AvatarURL:   &avatarUrls[i],
				Address:     &addresses[i],
				Birthday:    &birthdays[i],
				PhoneNumber: &phoneNumbers[i],
				RoleName:    roleNames[i],
			})
	}
	SeedAccounts(db, &accounts)
	SeedRolesForAccounts(db, &accounts)
	SeedPricings(db)
	SeedCoinBundles(db)
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
	fmt.Println("======= Role seeded.")
}

func SeedPricings(db *gorm.DB) {

	pricings := []models.Pricing{

		{
			ServiceName:  "messaging_session",
			Quantity:     1,
			QuantityUnit: "session",
			Price:        30,
			PriceUnit:    "coin(s)",
		},
		{
			ServiceName:  "translation_session",
			Quantity:     1,
			QuantityUnit: "session",
			Price:        50,
			PriceUnit:    "coin(s)",
		},
		{
			ServiceName:  "private_call_session",
			Quantity:     1,
			QuantityUnit: "session",
			Price:        40,
			PriceUnit:    "coin(s)",
		},
	}
	db.Create(&pricings)
	fmt.Println("======= Pricings seeded.")
}

func SeedCoinBundles(db *gorm.DB) {

	bundles := []models.CoinBundle{

		{
			Title:       "Pack of Coins",
			Description: "Recommended for new users to try out.",
			Quantity:    10,
			Price:       10000,
			PriceUnit:   "VND",
		},
		{
			Title:       "Pocket of Coins",
			Description: "Recommended for active users.",
			Quantity:    20,
			Price:       15000,
			PriceUnit:   "VND",
		},
		{
			Title:       "Bunch of Coins",
			Description: "Best offer available at the moment.",
			Quantity:    60,
			Price:       30000,
			PriceUnit:   "VND",
		},
	}
	db.Create(&bundles)
	fmt.Println("======= Coin bundles seeded.")
}
