package author

import (
	"fiberStore/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func InitMySQL() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err.Error())
	}

	config := Config{
		DB_Username: os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_Host:     os.Getenv("DB_HOST"),
		DB_Name:     os.Getenv("DB_NAME"),
	}

	ConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)

	dbConn, err := gorm.Open(mysql.Open(ConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	DB = dbConn.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return time.Now().In(location)
		},
	})

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Error getting underlying database connection:", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	return DB, nil
}

func AccountSeeder(db *gorm.DB) error {
	password := "$2a$10$RQ.ycP2DqXWYYDdcYuO6hebdgWOTIPpMEE8d6.ntVTDE21f6RkZKG"
	users := []models.User{
		{
			Name:     "Administrator",
			Username: "admin",
			Password: password,
			Role:     "Admin",
		},
		{
			Name:     "Customer",
			Username: "customer",
			Password: password,
			Role:     "Customer",
		},
	}

	for _, user := range users {
		// Check Data if already seeding
		var count int64
		if err := db.Model(&models.User{}).Where(&user).Count(&count).Error; err != nil {
			return err
		}

		// If data exists, skip seeding
		if count > 0 {
			continue
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}

		// Create User Amount
		userAmount := models.UserAmount{
			UserID: user.ID,
			Amount: 0,
		}

		// Check user amount if already seeding
		var countUserAmount int64
		if err := db.Model(&models.UserAmount{}).Where(&userAmount).Count(&countUserAmount).Error; err != nil {
			return err
		}

		// If data exists, skip seeding
		if countUserAmount > 0 {
			continue
		}

		if err := db.Create(&userAmount).Error; err != nil {
			return err
		}

	}

	return nil
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.UserAmount{},
		&models.Product{},
		&models.Cart{},
		&models.CartDetail{},
		&models.Transaction{},
		&models.TransactionDetail{},
	)
}
