package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  	string `json:"username" binding:"required" gorm:"unique"`
	Password	string
	Firstname	string
	Lastname 	string
	BankAccountNo string `gorm:"unique"`
	CreditBalance float32
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


type UserUpdate struct {
	Password string `json:"password" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname string `json:"lastname" binding:"required"`
	BankAccountNo string `json:"bankaccountno" binding:"required"`
}