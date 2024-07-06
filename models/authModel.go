package models

// User represents the user model
type UserAuth struct {
	UserId   string `json:"userId" gorm:"column:userId;primaryKey"`
	Email    string `json:"email" gorm:"column:email" validate:"required"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Provider string `json:"provider" gorm:"column:provider"`
}

// TableName sets the insert table name for this struct type
func (UserAuth) TableName() string {
	return "users"
}

type UserLogin struct {
	Username string `json:"username" gorm:"column:username" validate:"required"`
	Password string `json:"password" gorm:"column:password" validate:"required"`
}
