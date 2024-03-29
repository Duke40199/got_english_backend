package database

import (
	"fmt"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/models"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

var roleNameConfig = config.GetRoleNameConfig()
var priceEachCoin uint = 2000

//SeedDB function will trigger all seed functions below
func SeedDB(db *gorm.DB) {
	accounts := []models.Account{}
	uuid1, _ := uuid.Parse("1f78cabc-b268-43cb-9935-c3a0a53f4f82")
	uuid2, _ := uuid.Parse("0edb6398-fa61-43c9-9ffd-e83127fc6060")
	uuid3, _ := uuid.Parse("cd9d9123-a7cc-48ed-87e1-045b21eaf466")
	uuid4, _ := uuid.Parse("4b481c87-f208-4dfe-bc44-18c631a95a34")
	ids := []uuid.UUID{uuid1, uuid2, uuid3, uuid4}
	usernames := []string{"TuanAnh", "DucPhi", "TuanNguyen", "LocTr"}
	fullnames := []string{"Nguyễn Tuấn Anh", "Phí Đỗ Hồng Đức", "Nguyễn Trần Quốc Tuấn", "Trần Thiên Lộc"}
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
	birthdays := []string{
		"1987-01-12",
		"1998-02-23",
		"1999-03-14",
		"1991-10-12",
	}
	for i := 0; i < len(usernames); i++ {
		accounts = append(accounts,
			models.Account{
				ID:          ids[i],
				Username:    &usernames[i],
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
	SeedExchangeRates(db)
	SeedRatingAlgorithm(db)
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
					AccountID:          (*accounts)[i].ID,
					CanManageExpert:    true,
					CanManageLearner:   true,
					CanManageAdmin:     true,
					CanManageModerator: true,
				})
				break
			}
		case roleNameConfig.Learner:
			{
				db.Create(&models.Learner{
					AccountID:          (*accounts)[i].ID,
					AvailableCoinCount: 10000,
				})
				break
			}
		case roleNameConfig.Expert:
			{
				db.Create(&models.Expert{
					AccountID:                 (*accounts)[i].ID,
					CanChat:                   true,
					CanJoinTranslationSession: true,
					CanJoinLiveCallSession:    true,
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
					CanManageExchangeRate:    true,
					CanManageRatingAlgorithm: true,
				})
				break
			}
		}
	}
	fmt.Println("======= Role seeded.")
}

func SeedPricings(db *gorm.DB) {

	pricings := []models.Pricing{
		//Messaging
		{
			PricingName:  "messaging_session",
			Quantity:     1,
			QuantityUnit: "session",
			Price:        30,
			PriceUnit:    "coin(s)",
		},
		//LiveCall
		{
			PricingName:  "live_call_session",
			Quantity:     5,
			QuantityUnit: "minutes",
			Price:        40,
			PriceUnit:    "coin(s)",
		},
		{
			PricingName:  "live_call_session",
			Quantity:     10,
			QuantityUnit: "minutes",
			Price:        50,
			PriceUnit:    "coin(s)",
		},
		{
			PricingName:  "live_call_session",
			Quantity:     30,
			QuantityUnit: "minutes",
			Price:        60,
			PriceUnit:    "coin(s)",
		},
		//Translation
		{
			PricingName:  "translation_call_session",
			Quantity:     5,
			QuantityUnit: "minutes",
			Price:        60,
			PriceUnit:    "coin(s)",
		},
		{
			PricingName:  "translation_call_session",
			Quantity:     10,
			QuantityUnit: "minutes",
			Price:        70,
			PriceUnit:    "coin(s)",
		},
		{
			PricingName:  "translation_call_session",
			Quantity:     30,
			QuantityUnit: "minutes",
			Price:        80,
			PriceUnit:    "coin(s)",
		},
		{
			PricingName:  "coin_value",
			Quantity:     1,
			QuantityUnit: "coin",
			Price:        uint(priceEachCoin),
			PriceUnit:    "VND",
		},
	}
	db.Create(&pricings)
	fmt.Println("======= Pricings seeded.")
}

func SeedRatingAlgorithm(db *gorm.DB) {
	var minimumRatingCount uint = 5
	var averageAllExpertsRating float32 = 0
	algorithm := models.RatingAlgorithm{
		ID:                      1,
		MinimumRatingCount:      minimumRatingCount,
		AverageAllExpertsRating: averageAllExpertsRating,
	}
	db.Create(&algorithm)
	fmt.Println("======= Rating algorithm seeded.")
}

func SeedCoinBundles(db *gorm.DB) {
	titles := []string{"Gói Khởi Động", "Gói Chuyên Nghiệp", "Gói Lâu Dài"}
	descriptions := []string{"Phù hợp với những người dùng mới bắt đầu", "Phù hợp với những người dùng tần suất cao", "Phù hợp với những người dùng lâu dài"}
	quantities := []uint{10, 20, 30}
	pricings := []uint{10 * priceEachCoin, 20 * priceEachCoin, 30 * priceEachCoin}
	priceUnit := "VND"
	bundles := []models.CoinBundle{

		{
			Title:       &titles[0],
			Description: &descriptions[0],
			Quantity:    &quantities[0],
			Price:       &pricings[0],
			PriceUnit:   &priceUnit,
		},
		{
			Title:       &titles[1],
			Description: &descriptions[1],
			Quantity:    &quantities[1],
			Price:       &pricings[1],
			PriceUnit:   &priceUnit,
		},
		{
			Title:       &titles[2],
			Description: &descriptions[2],
			Quantity:    &quantities[2],
			Price:       &pricings[2],
			PriceUnit:   &priceUnit,
		},
	}
	db.Create(&bundles)
	fmt.Println("======= Coin bundles seeded.")
}

func SeedExchangeRates(db *gorm.DB) {
	rateList := []float32{0.2, 0.3, 0.4}
	rates := []models.ExchangeRate{}
	serviceNameList := []string{"messaging_session", "live_call_session", "translation_session"}
	for i := 0; i < len(rateList); i++ {
		rates = append(rates, models.ExchangeRate{
			Rate:        &rateList[i],
			ServiceName: &serviceNameList[i],
		})
	}
	db.Create(&rates)
	fmt.Println("======= Exchange rates seeded.")
}
