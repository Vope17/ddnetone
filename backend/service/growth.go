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

type MilestoneResult struct {
	Target    int    `json:"target"`
	Timestamp string `json:"timestamp"`
	Maps      int    `json:"maps"`
}

func GetMilestones(c *gin.Context) {
	const TARGET_MAPS = 2403
	results := []MilestoneResult{}

	for target := 100; target <= TARGET_MAPS; target += 100 {
		var record model.GrowthData
		err := db.GetDB().
			Where("maps >= ?", target).
			Order("maps asc, id asc").
			First(&record).Error
		if err == nil {
			results = append(results, MilestoneResult{
				Target:    target,
				Timestamp: record.Timestamp,
				Maps:      record.Maps,
			})
		}
	}

	c.JSON(http.StatusOK, results)
}

func GetScoreMilestones(c *gin.Context) {
	const SCORE_STEP = 1000
	const SCORE_MAX = 100000
	results := []MilestoneResult{}

	for target := SCORE_STEP; target <= SCORE_MAX; target += SCORE_STEP {
		var record model.GrowthData
		err := db.GetDB().
			Where("points >= ?", target).
			Order("points asc, id asc").
			First(&record).Error
		if err == nil {
			results = append(results, MilestoneResult{
				Target:    target,
				Timestamp: record.Timestamp,
				Maps:      record.Points,
			})
		} else {
			break // 超過目前最高分，停止
		}
	}

	c.JSON(http.StatusOK, results)
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
