package database

import (
	"example/backend/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database connection
var db *gorm.DB

// InitDB initializes the database connection
func InitMySQL() {
	var err error

	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	credentials := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
	db, err = gorm.Open(mysql.Open(credentials), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// Auto migrate the models (add more models here as needed)
	db.AutoMigrate(&models.User{})
}

// GetDB returns the database connection
func MysqlDB() *gorm.DB {
	return db
}
