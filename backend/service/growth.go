package service

import (
	"net/http"
	"time"

	"DDNETONE/db"
	"DDNETONE/model"
	"github.com/gin-gonic/gin"
)

func GetGrowth(c *gin.Context) {
	var growth []model.GrowthData
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	db.GetDB().Where("timestamp >= ?", sevenDaysAgo).Order("id asc").Find(&growth)
	c.JSON(http.StatusOK, growth)
}

func RecordGrowthSnapshot(score int, maps int, runner string, map_name string, map_points int) {
	database := db.GetDB()
	var firstRecord model.MapRecord
	var startTime time.Time

	err := database.Where("status = 2 AND finish_time IS NOT NULL").Order("finish_time asc").First(&firstRecord).Error
	if err == nil && firstRecord.FinishTime != nil {
		startTime = *firstRecord.FinishTime
	} else {
		startTime = time.Now()
	}

	hoursSinceStart := time.Since(startTime).Hours()
	if hoursSinceStart < 0 {
		hoursSinceStart = 0
	}

	var lastGrowth model.GrowthData
	if err := database.Order("id desc").First(&lastGrowth).Error; err == nil {
		if lastGrowth.Points == score && lastGrowth.Maps == maps {
			return
		}
	}

	newGrowth := model.GrowthData{
		Hours:     hoursSinceStart,
		Points:    score,
		Maps:      maps,
		Runner:    runner,
		MapName:   map_name,
		MapPoints: map_points,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	database.Create(&newGrowth)
}
