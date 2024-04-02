package models

import (
	"fmt"

	"example/web-service-gin/utils/token"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func (u *User) SaveUser() (*User, error) {

	var err error
	fmt.Println("Error creating user:", DB) // Log the error for debugging

	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}
func LoginCheck(username string, password string) (string, error) {

	var err error

	u := User{}
	// Check if the user with given username exists
	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}
	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// Passwords do not match
		return "", err
	}
	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}
func GetUserByID(uid uint) (User, error) {

	var u User
	// Get user data with given id
	if err := DB.Select("id, username, created_at, updated_at").First(&u, uid).Error; err != nil {
		return u, err
	}
	return u, nil

}
