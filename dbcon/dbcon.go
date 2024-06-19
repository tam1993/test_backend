package dbcon

import (
	_ "fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("can't load .env file")
	}
	// connect db
	var dbName = os.Getenv("DB_NAME")
	var dbUserName = os.Getenv("DB_USERNAME")
	var dbPassword = os.Getenv("DB_PASSWORD")
	var dbHost = os.Getenv("DB_HOST")

	dsn := dbUserName + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName + "?charset=utf8mb4&parseTime=True"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connection failed")
		return false
	}
	Db.AutoMigrate()
	return true
}

func Migrate() {
	Db.AutoMigrate()
}
