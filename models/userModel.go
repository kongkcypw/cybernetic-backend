package models

// User represents the user model
type User struct {
	UserId      string `json:"userId" gorm:"column:userId;primaryKey"`
	FirstName   string `json:"firstName" gorm:"column:firstName" validate:"required"`
	LastName    string `json:"lastName" gorm:"column:lastName" validate:"required"`
	Email       string `json:"email" gorm:"column:email" validate:"required"`
	Password    string `json:"password" gorm:"column:password" validate:"required"`
	PhoneNumber string `json:"phoneNumber" gorm:"column:phoneNumber" validate:"required"`
}
