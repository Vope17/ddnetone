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
	var loadedScore int64
	var completedCount int64
	var loadedCount int64

	database := db.GetDB()

	database.Model(&model.MapRecord{}).Where("status = 2").Select("COALESCE(SUM(points), 0)").Scan(&completedScore)
	database.Model(&model.MapRecord{}).Where("status = 2").Count(&completedCount)
	database.Model(&model.MapRecord{}).Where("status = 3").Select("COALESCE(SUM(points), 0)").Scan(&loadedScore)
	database.Model(&model.MapRecord{}).Where("status = 3").Count(&loadedCount)

	var summary model.Summary
	if err := database.Last(&summary).Error; err != nil {
		summary = model.Summary{
			TargetScore: 10000,
		}
		database.Create(&summary)
	}

	// Phase 2 (加載階段): 有任何已加載地圖時，從 10000 倒數
	if loadedScore > 0 {
		summary.CurrentScore = 10000 - int(loadedScore)
	} else {
		summary.CurrentScore = int(completedScore)
	}

	summary.CompletedMaps = int(completedCount)
	summary.LoadedMaps = int(loadedCount)
	summary.LastUpdate = time.Now()

	database.Save(&summary)

}
