package model

type User struct {
	ID           int64  `gorm:"column:id"`
	Login        string `gorm:"column:login"`
	PasswordHash string `gorm:"column:pass_hash"`
	Email        string `gorm:"column:email"`
}
