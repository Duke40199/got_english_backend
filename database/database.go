package database

import (
	"fmt"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//ConnectToDB will connect the BE to the database
func ConnectToDB() (*gorm.DB, error) {
	config := config.GetConfig()
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DatabaseUsername,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		fmt.Println("DB connected.")
		// fmt.Printf("DBURI: %s", dsn)
		return db, nil
	} else {
		fmt.Println("ERROR FOUND!")
		panic(err)
	}
}

//SyncDB will migrate & seed DB
func SyncDB(isForced bool) {
	modelList := models.GetModelList()
	db, err := ConnectToDB()
	if isForced {

		//DB will be synced forcefully.
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(modelList); i++ {
			if db.Migrator().HasTable(modelList[i]) {
				db.Migrator().DropTable(modelList[i])
				db.Set("gorm:table_options", " DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci").Migrator().AutoMigrate(modelList[i])
			} else {
				db.Set("gorm:table_options", " DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci").Migrator().CreateTable(modelList[i])
			}
		}
		db.Set("gorm:table_options", " DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci").Migrator().AutoMigrate(&models.TranslationSession{})
		if !db.Migrator().HasTable(&models.Earning{}) {
			db.Set("gorm:table_options", " DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci").Migrator().CreateTable(&models.Earning{})
		}
		SeedDB(db)
		//Setup custom relations
		// db.Exec("ALTER TABLE `experts` ADD CONSTRAINT `fk_expert_accounts1` FOREIGN KEY (`accounts_id`) REFERENCES `accounts` (`id`)")
		// db.Exec("ALTER TABLE `learners` ADD CONSTRAINT `fk_learner_accounts1` FOREIGN KEY (`accounts_id`) REFERENCES `accounts` (`id`)")
		// db.Exec("ALTER TABLE `moderators` ADD CONSTRAINT `fk_moderator_accounts1` FOREIGN KEY (`accounts_id`) REFERENCES `accounts` (`id`)")
		// db.Exec("ALTER TABLE `admins` ADD CONSTRAINT `fk_admin_accounts1` FOREIGN KEY (`accounts_id`) REFERENCES `accounts` (`id`)")

	} else {
		fmt.Println("No seeding needed.")
		// fmt.Println("seeding...")
		// for i := 0; i < len(modelList); i++ {
		// 	if db.Migrator().HasTable(modelList[i]) {
		// 		db.Migrator().CreateTable(modelList[i])
		// 	}
		// }
		//Relations raw queries
	}

}
