package db

import (
	"log"
	"project-kunduz/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=db user=user password=password dbname=kunduz port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panic("Failed to connect to database!")
	}

	database.AutoMigrate(&models.User{}, &models.Cluster{}, &models.SBOM{}, &models.Asset{}, &models.CVE{}, &models.Alarm{})

	DB = database
}
