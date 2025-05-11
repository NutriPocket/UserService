package test

import (
	"os"

	"github.com/NutriPocket/UserService/database"
	"github.com/joho/godotenv"
	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var log = logging.MustGetLogger("log")
var gormDB *gorm.DB

func loadEnv() {
	if ci_test := os.Getenv("CI_TEST"); ci_test != "" {
		return
	}

	err := godotenv.Load("../../../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func setupDB() {
	database.ConnectDB()
	var err error
	gormDB, err = database.GetPoolConnection()
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}
}

func ClearUsers() {
	if err := gormDB.Exec(`
		DELETE FROM users
	`).Error; err != nil {
		log.Fatal(err)
	}
}

func ClearBlacklist() {
	if err := gormDB.Exec(`
		DELETE FROM jwt_blacklist
	`).Error; err != nil {
		log.Fatal(err)
	}
}

func Setup(testType string) {
	log.Infof("Setup %s tests!\n", testType)
	loadEnv()
	log.Info(".env.test loaded")
	setupDB()
}

func TearDown(testType string) {
	log.Infof("Tear down %s tests!\n", testType)
	database.Close()
}
