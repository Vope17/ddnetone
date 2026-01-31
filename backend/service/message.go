package service

import (
	"net/http"
	"time"

	"DDNETONE/db"
	"DDNETONE/model"
	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {

	var messages []model.Message
	db.GetDB().Order("created_at desc").Find(&messages)
	c.JSON(http.StatusOK, messages)
}

func CreateMessage(c *gin.Context) {
	var msg model.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg.CreatedAt = time.Now()
	if err := db.GetDB().Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post message"})
		return
	}
	c.JSON(http.StatusCreated, msg)
}
