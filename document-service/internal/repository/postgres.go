package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BHAV0207/documet-service/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dns := os.Getenv("URI")
	if dns == "" {
		dns = "host=aws-1-ap-southeast-1.pooler.supabase.com user=postgres.tnphmjbddmqgpesqrtgp password=House@12345 dbname=postgres port=5432 sslmode=require"

	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// Auto migrate
	if err := db.AutoMigrate(&models.Document{}, &models.Snapshot{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	DB = db
	fmt.Println("Database connected")

}
