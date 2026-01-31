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
	var totalScore int64

	var completedCount int64

	database := db.GetDB()

	database.Model(&model.MapRecord{}).Where("status = 2").Select("COALESCE(SUM(points), 0)").Scan(&totalScore)
	database.Model(&model.MapRecord{}).Where("status = 2").Count(&completedCount)

	var summary model.Summary
	if err := database.Last(&summary).Error; err != nil {
		summary = model.Summary{
			TargetScore: 10000,
		}
		database.Create(&summary)
	}

	summary.CurrentScore = int(totalScore)
	summary.CompletedMaps = int(completedCount)
	summary.LastUpdate = time.Now()

	database.Save(&summary)

}
