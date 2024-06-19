package main


import (
	"github.com/tam1993/test_backend/controller/transaction"
	"github.com/tam1993/test_backend/controller/user"
	"github.com/tam1993/test_backend/dbcon"
	"github.com/tam1993/test_backend/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var saltPass = "saltttttP"

func main() {
	setDb := dbcon.InitDB()
	if !setDb {
		fmt.Println("DB connection fail")
		return
	}
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//migrate
	r.GET("/migrate", func(c *gin.Context) {
		dbcon.Db.AutoMigrate(&model.User{})
		dbcon.Db.AutoMigrate(&model.UserTopupTransaction{})
		dbcon.Db.AutoMigrate(&model.UserTransferCreditTransaction{})
	})

	r.POST("/user/register", user.Register)
	r.POST("/user/login", user.Login)

	authUser := r.Group("/user", Auth())
	authUser.GET("/me", user.Me)
	authUser.PATCH("/me", user.UpdateUser)

	authAccounting := r.Group("/accounting", Auth())
	authAccounting.POST("/topup", transaction.Topup)
	authAccounting.POST("/transfer", transaction.Transfer)
	authAccounting.GET("/transfer-list", transaction.TransferList)
	// // r.POST("/transaction/usertransaction/", transaction.SelectUserTransaction)
	// auth.POST("/add", transaction.AddTransaction)
	// auth.POST("/usertransaction", transaction.SelectUserTransaction)


	r.Run(":8888")
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//check auth
		tokenString := strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", -1)
		hmacSampleSecret := []byte(saltPass)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return hmacSampleSecret, nil
		})
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		c.Set("userID", claims["userID"])
		c.Next()
	}
}