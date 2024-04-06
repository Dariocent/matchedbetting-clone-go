package initializers

import "github.com/Dariocent/matchedbetting-clone-go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
