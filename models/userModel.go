package models

// User represents the user model
type User struct {
	UserId   string `json:"userId" gorm:"column:userId;primaryKey"`
	Email    string `json:"email" gorm:"column:email" validate:"required"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Provider string `json:"provider" gorm:"column:provider"`
}

type UserLogin struct {
	Username string `json:"username" gorm:"column:username" validate:"required"`
	Password string `json:"password" gorm:"column:password" validate:"required"`
}
