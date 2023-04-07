package controllers

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"lmcapi/models"
	L "math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CheckTokenInput struct {
	Token string `json:"token" binding:"required"`
}

type CheckUuidInput struct {
	Uuid  string `json:"uuid" binding:"required"`
	Token string `json:"token" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// AllUsers GET /users
// Find all users
func AllUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)
	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	for i, user := range users {
		user.Password = ""
		users[i] = user
	}
	c.JSON(http.StatusOK, gin.H{"users": users, "success": true})
}

// Login POST /user/login
// Login user
func Login(c *gin.Context) {
	ip := c.Request.Header.Get("X-Forwarded-For")
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	var user models.User
	models.DB.Where("email = ? AND password = ?", input.Email, makeHash(input.Password)).First(&user)
	if (models.User{} == user) {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	token := makeToken(ip)
	token.User = user
	if err := models.DB.Create(&token).Error; err != nil {
		if strings.Contains(err.Error(), "duplicated key not allowed") {
			var tokenDb models.Token
			models.DB.Where("user_id = ?", user.ID).First(&tokenDb)
			models.DB.Delete(&tokenDb)
			models.DB.Create(&token)
		}
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token.Token, "success": true})
}

// Register POST /user/register
// Register user
func Register(c *gin.Context) {
	ip := c.Request.Header.Get("X-Forwarded-For")
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	var user models.User
	user.Name = input.Name
	user.Email = input.Email
	user.Password = makeHash(input.Password)
	user.Balance = 0.0
	user.SupportCode = makeKey()
	if err := models.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicated key not allowed") {
			c.JSON(http.StatusOK, gin.H{"success": false})
			return
		}
	}
	models.DB.Where("email = ?", input.Email).First(&user)
	token := makeToken(ip)
	token.User = user
	if err := models.DB.Create(&token).Error; err != nil {
		if strings.Contains(err.Error(), "duplicated key not allowed") {
			var tokenDb models.Token
			models.DB.Where("user_id = ?", user.ID).First(&tokenDb)
			models.DB.Delete(&tokenDb)
			models.DB.Create(&token)
		}
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token.Token, "success": true})
}

// CheckToken POST /users/checktoken
// Checks given Token
func CheckToken(c *gin.Context) {
	ip := c.Request.Header.Get("X-Forwarded-For")
	var input CheckTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	if !compareToken(input.Token, ip) {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	var tokenDb models.Token
	models.DB.Preload("User").Where("token = ?", input.Token).First(&tokenDb)
	user := tokenDb.User
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"user": user, "token": input.Token, "success": true})
}

// GetUserDataByToken POST /user/getuserdatabytoken
// Gets User Data by Token
func GetUserDataByToken(c *gin.Context) {
	ip := c.Request.Header.Get("X-Forwarded-For")
	var input CheckTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	if !compareToken(input.Token, ip) {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	var tokenDb models.Token
	models.DB.Preload("User").Where("token = ?", input.Token).First(&tokenDb)
	user := tokenDb.User
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"user": user, "success": true})
}

// makeHash
// Make SHA256 from .env SALT and given password
func makeHash(password string) string {
	h := sha256.New()
	h.Write([]byte(os.Getenv("SALT") + password))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

// makeToken
// Make SHA256 from .env SALT ip and decoderkey
func makeToken(ip string) models.Token {
	b := make([]byte, 5)
	_, _ = rand.Read(b)
	decoderKey := hex.EncodeToString(b)
	h := sha256.New()
	h.Write([]byte(os.Getenv("SALT") + ip + decoderKey))
	tk := hex.EncodeToString(h.Sum(nil))
	token := models.Token{Token: tk, DecoderKey: decoderKey}
	return token
}

// compareToken
// Compare given Token
func compareToken(token string, ip string) bool {
	var tokenDb models.Token
	models.DB.Where("token = ?", token).First(&tokenDb)
	if (models.Token{} == tokenDb) {
		return false
	}
	h := sha256.New()
	h.Write([]byte(os.Getenv("SALT") + ip + tokenDb.DecoderKey))
	tk := hex.EncodeToString(h.Sum(nil))
	return token == tk
}

func makeKey() string {
	r := L.New(L.NewSource(time.Now().UnixNano()))
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	parts := make([][]byte, 3)
	for i := range parts {
		part := make([]byte, 4)
		for j := range part {
			part[j] = chars[r.Intn(len(chars))]
		}
		parts[i] = part
	}
	return string(bytes.Join(parts, []byte("-")))
}
