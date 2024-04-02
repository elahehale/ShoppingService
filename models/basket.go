package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type Basket struct {
	gorm.Model
	UserID uint        `json:"user_id" validate:"required" gorm:"foreignkey:UserID"`
	Data   []byte      `json:"data" gorm:"type:jsonb"`
	State  BasketState `json:"state,omitempty" gorm:"default:'PENDING'"`
}
type BasketJSON struct {
	ID     uint        `json:"id"`
	UserID uint        `json:"user_id"`
	Data   interface{} `json:"data"`
	State  BasketState `json:"state"`
}
type BasketState string
type JSONB map[string]interface{}

const (
	Completed BasketState = "COMPLETED"
	Pending   BasketState = "PENDING"
)

// Function for serializing basket data field
func (b *Basket) BasketSerializer() (*BasketJSON, error) {
	var jsonData interface{}
	err := json.Unmarshal(b.Data, &jsonData)
	if err != nil {
		return nil, err
	}

	basketJSON := &BasketJSON{
		ID:     b.ID,
		UserID: b.UserID,
		Data:   jsonData,
		State:  b.State,
	}
	return basketJSON, nil
}

func (b *Basket) SaveBasket() (*BasketJSON, error) {
	var err error
	fmt.Println("Error creating user:", DB)

	if b.ID != 0 {
		// Basket exists, update it
		err = DB.Save(b).Error
	} else {
		// Basket does not exist, create a new one
		err = DB.Create(b).Error
	}
	if err != nil {
		return nil, err
	}

	basketJSON, err := b.BasketSerializer()
	if err != nil {
		return nil, err
	}
	return basketJSON, nil
}

func GetBasketByID(uid uint) (*BasketJSON, error) {

	var b Basket

	if err := DB.First(&b, uid).Error; err != nil {
		return nil, errors.New("Basket not found!")
	}
	// serialize data of basket
	basketJSON, err := b.BasketSerializer()
	if err != nil {
		return nil, err
	}
	return basketJSON, nil

}
func GetBasketObjectByID(uid uint) (*Basket, error) {

	var b Basket
	// Check if the basket exists
	if err := DB.First(&b, uid).Error; err != nil {
		return nil, errors.New("Basket not found!")
	}

	return &b, nil

}
func DeleteBasket(basketID uint) error {
	var b Basket

	// Check if the basket exists
	if err := DB.First(&b, basketID).Error; err != nil {
		return errors.New("Basket not found")
	}

	// Delete the basket from the database
	if err := DB.Delete(&b).Error; err != nil {
		return errors.New("Failed to delete basket")
	}

	return nil
}

func GetBasketsByUserID(userID int64) ([]BasketJSON, error) {
	var baskets []Basket

	DB = DB.Debug()

	// Query the database to get baskets with the given user ID
	if err := DB.Select("id, created_at, updated_at, user_id, data, state").Where("user_id = ?", userID).Find(&baskets).Error; err != nil {
		return nil, fmt.Errorf("Failed to retrieve baskets: %v", err)
	}

	// Convert each Basket to BasketJSON
	var basketJSONList []BasketJSON
	for _, b := range baskets {
		basketJSON, err := b.BasketSerializer()
		if err != nil {
			return nil, err
		}
		basketJSONList = append(basketJSONList, *basketJSON)
	}

	return basketJSONList, nil
}
