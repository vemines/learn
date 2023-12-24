package initializers

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnnectToDb() {
	dsn := os.Getenv("DNS")
	var err error
	// Create a new GORM database connection.
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed connect to MySQL:")
	}
}
