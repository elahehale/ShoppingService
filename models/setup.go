package models

import (
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq"
)

// postgres database configurations
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "webdb"
)

var DB *gorm.DB

func ConnectDataBase() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the database
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	DB = db
	fmt.Println("Successfully connected to postgres!")

	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Basket{})
	DB.Model(&Basket{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	DB.Exec("ALTER TABLE baskets ADD CONSTRAINT check_data_size CHECK (LENGTH(CAST(data AS text)) <= 2048)")

	return DB
}
