package user

import (
	"github.com/tam1993/test_backend/dbcon"
	"github.com/tam1993/test_backend/model"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// var saltPass = os.Getenv("PASSWORD_SALT")
var saltPass = "saltttttP"

func Register(c *gin.Context) {
	var json model.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Username
	isValidUsername := regexp.MustCompile(`^[a-z0-9]+$`).MatchString(json.Username)
	if !isValidUsername {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be a-z 0-9"})
		return
	}

	// Validate BankAccountNo
	isValidBankAccountNo := regexp.MustCompile(`^\d{10}$`).MatchString(json.BankAccountNo)
	if !isValidBankAccountNo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BankAccountNo must be 0-9"})
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(saltPass + json.Password))
	encryptPassWord := hex.EncodeToString(hasher.Sum(nil))

	UserData := model.User{
		Username: json.Username,
		Password: encryptPassWord,
		Firstname: json.Firstname,
		Lastname: json.Firstname,
		BankAccountNo: json.BankAccountNo,
		//credit free 1000
		CreditBalance: 1000,
	}
	dbcon.Db.Create(&UserData)
	if UserData.ID == 0 {
		c.JSON(400, gin.H{"message": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func Login(c *gin.Context) {
	var json model.UserLogin
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Username
	isValidUsername := regexp.MustCompile(`^[a-z0-9]+$`).MatchString(json.Username)
	if !isValidUsername {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be a-z 0-9"})
		return
	}

	var userDetail model.User
	result := dbcon.Db.Where("username = ?", json.Username).First(&userDetail)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password incorrect"})
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(saltPass + json.Password))
	encryptPassWord := hex.EncodeToString(hasher.Sum(nil))
	if encryptPassWord != userDetail.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password incorrect"})
		return
	}

	hmacSampleSecret := []byte(saltPass)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userDetail.ID,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(200, gin.H{"token": tokenString})
}

func Me(c *gin.Context) {
	userId := int64(c.MustGet("userID").(float64))

	var users model.User
	result := dbcon.Db.First(&users, userId)
	if result.Error != nil {
		c.JSON(400, result.Error)
		return
	}
	c.JSON(200, users)
}

func UpdateUser(c *gin.Context) {
	var json model.UserUpdate
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Validate BankAccountNo
	isValidBankAccountNo := regexp.MustCompile(`^\d{10}$`).MatchString(json.BankAccountNo)
	if !isValidBankAccountNo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BankAccountNo must be 0-9"})
		return
	}

	//check user is exists?
	var findUser model.User
	userId := int64(c.MustGet("userID").(float64))
	checkUser := dbcon.Db.First(&findUser, userId)
	if checkUser.Error != nil {
		c.JSON(400, gin.H{"Error": "not found"})
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(saltPass + json.Password))
	encryptPassWord := hex.EncodeToString(hasher.Sum(nil))

	findUser.Password = encryptPassWord
	findUser.Firstname = json.Firstname
	findUser.Lastname = json.Lastname
	findUser.BankAccountNo = json.BankAccountNo
	dbcon.Db.Save(findUser)
	c.JSON(200, findUser)
}