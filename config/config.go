package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

//Config BE model
type Config struct {
	AppName     string
	AppVersion  string
	Environment string
	APIHost     string `toml:"api_host"`
	JWTSecret   string `toml:"jwt_secret"`

	//Configurations for DB
	DatabaseHost     string `toml:"ge_database_host"`
	DatabaseName     string `toml:"ge_database_name"`
	DatabaseUsername string `toml:"ge_database_user"`
	DatabasePort     string `toml:"ge_database_port"`
	DatabasePassword string `toml:"ge_database_password"`
	DatabaseSslMode  string `toml:"ge_database_ssl_mode"`
	DatabaseDialect  string `toml:"ge_database_dialect"`
	DatabaseTimezone string `toml:"ge_database_timezone"`
}

func (c *Config) Read() {
	var configFile string

	switch env := os.Getenv("GO_ENV"); env {
	case "staging":
		configFile = "config.staging.toml"
	case "production":
		configFile = "config.production.toml"
	case "qa":
		configFile = "config.qa.toml"
	default:
		configFile = "config.local.toml"
	}
	basePath := os.Getenv("GOPATH") + "/src/github.com/golang/got_english_backend/"
	configFile = basePath + configFile

	if _, err := toml.DecodeFile(configFile, &c); err != nil {
		log.Fatal(err)
	}
	log.Println("configuration environment: ", c.Environment)
}

var config = Config{}

func init() {
	log.Println("Reading configuration file")
	config.Read()
	log.Println("Read configuration file")
}

//GetConfig : export backend config
func GetConfig() *Config {
	return &config
}

//RoleNameConfig model
type RoleNameConfig struct {
	Admin     string
	Moderator string
	Expert    string
	Learner   string
}

var roleNameConfig = RoleNameConfig{
	Admin:     "Admin",
	Moderator: "Moderator",
	Expert:    "Expert",
	Learner:   "Learner",
}

//GetRoleNameConfig : export roleID config
func GetRoleNameConfig() *RoleNameConfig {
	return &roleNameConfig
}

//RoleNameConfig model
type ApplicationFormStatusConfig struct {
	Pending  string
	Approved string
	Declined string
}

var applicationFormStatusConfig = ApplicationFormStatusConfig{
	Pending:  "Pending",
	Approved: "Approved",
	Declined: "Declined",
}

//GetRoleNameConfig : export roleID config
func GetApplicationFormStatusConfig() *ApplicationFormStatusConfig {
	return &applicationFormStatusConfig
}
