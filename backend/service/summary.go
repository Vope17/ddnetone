package service

import (
	"net/http"
	"time"

	"DDNETONE/db"
	"DDNETONE/model"
	"github.com/gin-gonic/gin"
)

// GetSummary Handler
func GetSummary(c *gin.Context) {
	UpdateGlobalSummary() // 呼叫內部邏輯更新
	var summary model.Summary
	db.GetDB().Last(&summary)

	c.JSON(http.StatusOK, summary)
}

// UpdateGlobalSummary 更新全服總覽數據 (內部呼叫用)
func UpdateGlobalSummary() {
	var completedScore int64
	var completedCount int64
	var totalScore int64
	var totalMaps int64

	database := db.GetDB()

	database.Model(&model.MapRecord{}).Where("status = 2").Select("COALESCE(SUM(points), 0)").Scan(&completedScore)
	database.Model(&model.MapRecord{}).Where("status = 2").Count(&completedCount)
	database.Model(&model.MapRecord{}).Select("COALESCE(SUM(points), 0)").Scan(&totalScore)
	database.Model(&model.MapRecord{}).Count(&totalMaps)

	var summary model.Summary
	if err := database.Last(&summary).Error; err != nil {
		summary = model.Summary{
			TargetScore: int(totalScore),
			TargetMaps:  int(totalMaps),
		}
		database.Create(&summary)
	}

	summary.CurrentScore = int(completedScore)
	summary.CompletedMaps = int(completedCount)
	summary.TargetScore = int(totalScore)
	summary.TargetMaps = int(totalMaps)
	summary.LastUpdate = time.Now()

	database.Save(&summary)

}
