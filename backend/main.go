// backend/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ContainerStatus struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	IPAddress      string    `json:"ip_address"`
	PingDuration   int64     `json:"ping_duration"` 
	LastSuccessful time.Time `json:"last_successful"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL") 
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&ContainerStatus{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func main() {
	initDB()
	r := gin.Default()

	r.GET("/containers", func(c *gin.Context) {
		var statuses []ContainerStatus
		if err := db.Find(&statuses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, statuses)
	})

	r.POST("/containers", func(c *gin.Context) {
		var status ContainerStatus
		if err := c.ShouldBindJSON(&status); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		status.LastSuccessful = time.Now()
		if err := db.Create(&status).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, status)
	})

	r.Run(":8080")
}

