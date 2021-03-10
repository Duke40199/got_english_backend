package database

import (
	"fmt"

	"github.com/golang/GotEnglishBackend/Application/config"
	"github.com/golang/GotEnglishBackend/Application/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

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
	fmt.Printf("+======DBRI:%s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		fmt.Printf("DBURI: %s", dsn)
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
				db.Migrator().CreateTable(modelList[i])
			} else {
				db.Migrator().CreateTable(modelList[i])
			}
		}
		SeedDB(db)
		//Setup custom relations
		db.Exec("ALTER TABLE `experts` ADD CONSTRAINT `fk_expert_accounts1` FOREIGN KEY (`accounts_id`) REFERENCES `accounts` (`id`)")
		db.Exec("ALTER TABLE `learners` ADD CONSTRAINT `fk_learner_accounts1` FOREIGN KEY (`accounts_id`) REFERENCES `accounts` (`id`)")
		db.Exec("ALTER TABLE `moderators` ADD CONSTRAINT `fk_moderator_accounts1` FOREIGN KEY (`accounts_id`) REFERENCES `accounts` (`id`)")
		db.Exec("ALTER TABLE `admins` ADD CONSTRAINT `fk_admin_accounts1` FOREIGN KEY (`accounts_id`) REFERENCES `accounts` (`id`)")
		db.AutoMigrate(&models.TranslationSession{}, &models.Learner{})
	} else {
		fmt.Println("seeding...")
		for i := 0; i < len(modelList); i++ {
			if db.Migrator().HasTable(modelList[i]) {
				db.Migrator().CreateTable(modelList[i])
			}
		}
		//Relations raw queries
	}

}
