package model

type  UserModel struct{
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	Email string `gorm:"unique"`
	Password string `gorm:"not null"`
	Status string `gorm:"not null"`
}