package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserBody struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func GenToken(db *gorm.DB, email string) string {

	var user User
	privateKey := LoadPrivateKey()

	result := db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		fmt.Printf("User not found\n")

		return ""
	}

	payload := Payload{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}

	token, err := GenJWT(privateKey, payload)

	if err != nil {
		fmt.Print(err.Error())
	}

	return token

}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://rest-orpin.vercel.app"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	db := ConnectDB("app.db")

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"name": "Mahadi",
		})
	})

	r.POST("/api/auth/sign-up/email", func(ctx *gin.Context) {

		var user UserBody
		err := ctx.ShouldBindJSON(&user)

		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err.Error(),
			})
		}

		if user.Email == "" || user.Name == "" || user.Password == "" {
			ctx.JSON(400, gin.H{
				"message": "email, name, and password is required",
			})
			return
		}

		db.Create(&User{
			Name:     user.Name,
			Password: user.Password,
			Email:    user.Email,
		})

		ctx.JSON(200, gin.H{
			"data":    user,
			"token":   GenToken(db, user.Email),
			"message": "User Created",
		})
	})

	r.GET("/me", func(c *gin.Context) {
		// 1. Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// 2. Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		tokenString := parts[1]

		// 3. Validate token
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 4. Send user data from token
		c.JSON(http.StatusOK, gin.H{
			"id":    claims.Subject,
			"name":  claims.Name,
			"email": claims.Email,
			"role":  claims.Role,
		})
	})

	fmt.Println("Server is running at :8080")
	r.Run(":8080")

}
