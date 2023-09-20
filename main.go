package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	r := gin.New()

	dsn := "root:kong9231@tcp(127.0.0.1:3306)/movie?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	r.GET("/books", getBooks)
	r.POST("/books", createBookHandler)

	r.Run()
}

type Movie struct {
	ID     *string `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
}

func getBooks(c *gin.Context) {
	var response *[]Movie

	if err := db.Table("movie").Select("*").Find(&response).Error; err != nil {
		c.JSON(http.StatusOK, nil)
		return
	}
	c.JSON(http.StatusOK, response)
}

func createBookHandler(c *gin.Context) {
	var request Movie

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("movie").Create(&request).Error; err != nil {
			c.JSON(http.StatusConflict, err)
			return nil
		} else {
			c.JSON(http.StatusCreated, nil)
			return nil;
		}
	})

	c.JSON(http.StatusBadGateway, err)
}
