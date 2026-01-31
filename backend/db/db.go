package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// 引入 model 包以進行 Migrate
	"DDNETONE/model" // ★請將 DDNETONE 替換為你的 go.mod module 名稱
)

var DB *gorm.DB

func Init() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbTimeZone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbTimeZone)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自動遷移 Schema
	DB.AutoMigrate(
		&model.Summary{},
		&model.Player{},
		&model.MapRecord{},
		&model.GrowthData{},
		&model.Message{},
	)

	log.Println("Database connected and migrated.")
}

// GetDB 提供給其他 package 使用
func GetDB() *gorm.DB {
	return DB
}
