package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"lmcapi/controllers"
	"lmcapi/customs"
	"lmcapi/models"
	"log"
	"os"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("\nError loading .env file\nPlease copy the example.env to .env and change your variables")
	}
	if os.Getenv("SALT") == "" {
		randomBytes := make([]byte, 80)
		_, err := rand.Read(randomBytes)
		if err != nil {
			log.Fatal("\nCouldn't create a salt")
		}
		salt := hex.EncodeToString(randomBytes)
		filePath := ".env"
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("\nCouldn't open .env file")
		}
		defer file.Close()
		if _, err = file.WriteString(fmt.Sprintf("\nSALT=\"%s\"", salt)); err != nil {
			log.Fatal("\nCouldn't write salt to .env file")
		}
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal("\nError loading .env file\nPlease copy the example.env to .env and change your variables")
		}
	}
}

func main() {
	r := gin.Default()
	models.ConnectDatabase()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{os.Getenv("APPLICATION_URL")},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Origin"},
	}))
	r.GET("/users", controllers.AllUsers)
	r.POST("/user/login", controllers.Login)
	r.POST("/user/register", controllers.Register)
	r.POST("/user/checktoken", controllers.CheckToken)
	r.POST("/user/getuserdatabytoken", controllers.GetUserDataByToken)
	customs.RegisterCustomRoutes(r)
	err := r.Run("127.0.0.1:3000")
	if err != nil {
		panic(err)
	}
}
