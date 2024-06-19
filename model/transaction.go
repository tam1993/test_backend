package model

import (
	"gorm.io/gorm"
)

type UserTopupTransaction struct {
	gorm.Model
	//foreign key
	UserID int64 `json:"userid" binding:"required"`
	User   User  `gorm:"foreignKey:UserID"`

	Amount       float32 `json:"amount" binding:"required"`
	AmountBefore float32
	AmountAfter  float32
}

type TopupTransaction struct {
	Amount       float32 `json:"amount" binding:"required"`
	AmountBefore float32
	AmountAfter  float32
}


type UserTransferCreditTransaction struct {
	gorm.Model
	//foreign key
	FromUserID int64 `json:"from_userid" binding:"required"`
	FromUser   User  `gorm:"foreignKey:FromUserID"`

	ToUserID int64 `json:"to_userid" binding:"required"`
	ToUser   User  `gorm:"foreignKey:ToUserID"`

	Amount       float32 `json:"amount" binding:"required"`
}

type TransferCreditTransaction struct {
	BankAccountNo string `json:"bankaccountno" binding:"required"`
	Amount       float32 `json:"amount" binding:"required"`
}

type TransferCreditTransactionList struct {
	FromDate string `json:"FromDate" binding:"required"`
	ToDate   string `json:"ToDate" binding:"required"`
}