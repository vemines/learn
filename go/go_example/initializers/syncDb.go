package initializers

import "go_example/models"

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}
