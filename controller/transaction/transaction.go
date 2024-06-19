package transaction

import (
	"github.com/tam1993/test_backend/dbcon"
	"github.com/tam1993/test_backend/model"
	// "database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Topup(c *gin.Context) {
	userId := int64(c.MustGet("userID").(float64))
	var json model.TopupTransaction
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var users model.User
	findUser := dbcon.Db.First(&users, userId)
	if findUser.Error != nil {
		c.JSON(400, findUser.Error)
		return
	}

	var TransactionData = model.UserTopupTransaction{
		UserID:       userId,
		Amount:       json.Amount,
		AmountBefore: users.CreditBalance,
	}
	if TransactionData.Amount < 0 {
		c.JSON(400, gin.H{})
		return
	}
	increaseUserCredit(&users, json.Amount)
	
	TransactionData.AmountAfter = users.CreditBalance

	dbcon.Db.Create(&TransactionData)
	if TransactionData.ID == 0 {
		c.JSON(400, TransactionData)
		return
	}

	dbcon.Db.Save(users)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}


func Transfer(c *gin.Context) {
	userId := int64(c.MustGet("userID").(float64))
	var json model.TransferCreditTransaction
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if(json.Amount <= 0){
		c.JSON(http.StatusBadRequest, gin.H{"message": "Amount error"})
		return
	}

	var fromUser model.User
	findUser := dbcon.Db.First(&fromUser, userId)
	if findUser.Error != nil {
		c.JSON(http.StatusBadRequest, findUser.Error)
		return
	}
	
	var toUser model.User
	if err := dbcon.Db.Where("bank_account_no = ? AND id != ?", json.BankAccountNo, fromUser.ID).First(&toUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bank Account not found"})
		return
	}

	if(json.Amount > fromUser.CreditBalance){
		c.JSON(http.StatusBadRequest, gin.H{"message": "Credit not enough"})
		return
	}

	var TransactionData = model.UserTransferCreditTransaction{
		FromUserID:	int64(fromUser.ID),
		ToUserID:	int64(toUser.ID),
		Amount: json.Amount,
	}
	
	dbcon.Db.Create(&TransactionData)
	if TransactionData.ID == 0 {
		c.JSON(400, TransactionData)
		return
	}

	decreaseUserCredit(&fromUser, json.Amount)
	dbcon.Db.Save(fromUser)
	increaseUserCredit(&toUser, json.Amount)
	dbcon.Db.Save(toUser)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func TransferList(c *gin.Context) {
	userId := int64(c.MustGet("userID").(float64))
	var json model.TransferCreditTransactionList
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var Transaction []model.UserTransferCreditTransaction

	fromDate, err := time.Parse("2006-01-02", json.FromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid FromDate format, use YYYY-MM-DD"})
		return
	}
	toDate, err := time.Parse("2006-01-02", json.ToDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ToDate format, use YYYY-MM-DD"})
		return
	}

	dbcon.Db.Preload("FromUser").Preload("ToUser").
	Where("from_user_id = ? OR to_user_id = ?", userId, userId).
	Where("DATE(created_at) >= ?", fromDate.Format("2006-01-02")).
	Where("DATE(created_at) <= ?", toDate.Format("2006-01-02")).
	Find(&Transaction)

	c.JSON(200, Transaction)
}

func increaseUserCredit(u *model.User, amount float32) {
	u.CreditBalance += amount
}

func decreaseUserCredit(u *model.User, amount float32) {
	u.CreditBalance -= amount
}
