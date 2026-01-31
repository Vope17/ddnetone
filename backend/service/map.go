package service

import (
	"net/http"
	"time"

	"DDNETONE/db"
	"DDNETONE/model"
	"github.com/gin-gonic/gin"
)

func GetMaps(c *gin.Context) {
	difficulty := c.Query("difficulty")
	var maps []model.MapRecord
	query := db.GetDB().Model(&model.MapRecord{})

	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	query.Order("score desc").Find(&maps)
	c.JSON(http.StatusOK, maps)
}

func GetMapOptions(c *gin.Context) {
	difficulty := c.Query("difficulty")
	var maps []model.MapRecord
	db.GetDB().Where("difficulty = ? AND status != 2", difficulty).
		Order("map_name asc").
		Find(&maps)
	c.JSON(http.StatusOK, maps)
}

func CreateRecord(c *gin.Context) {
	var newRecord model.MapRecord
	if err := c.ShouldBindJSON(&newRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	database := db.GetDB()
	var existingRecord model.MapRecord

	result := database.Where("map_name = ? AND difficulty = ? AND status != 2", newRecord.MapName, newRecord.Difficulty).First(&existingRecord)

	if result.Error == nil {
		// --- A. Update ---
		existingRecord.Runner = newRecord.Runner
		existingRecord.Note = newRecord.Note
		existingRecord.FinishTime = &now
		existingRecord.Status = newRecord.Status
		existingRecord.HasDummy = newRecord.HasDummy
		existingRecord.Score = newRecord.Score

		if existingRecord.Status == 2 {
			UpdatePlayerStats(existingRecord.Runner, existingRecord.Score)
		} else {
			existingRecord.Score = 0
			existingRecord.HasDummy = false
		}

		if err := database.Save(&existingRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}

		UpdateGlobalSummary()
		// 呼叫 Growth Snapshot (需先取得最新 Summary 數據)
		triggerSnapshot(existingRecord.Runner, existingRecord.MapName, existingRecord.Score)

		c.JSON(http.StatusOK, existingRecord)

	} else {
		// --- B. Create ---
		if err := database.Create(&newRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
			return
		}

		if newRecord.Status == 2 {
			UpdatePlayerStats(newRecord.Runner, newRecord.Score)
		}

		if newRecord.Score > 0 || newRecord.Status == 2 {
			newRecord.FinishTime = &now
			// 因為 BeforeSave 可能只有在 Save 時觸發較完整，這裡手動補強或依賴 GORM hooks
		}

		UpdateGlobalSummary()
		triggerSnapshot(newRecord.Runner, newRecord.MapName, newRecord.Score)

		c.JSON(http.StatusCreated, newRecord)
	}
}

// 輔助函式：取得當前數據並觸發快照
func triggerSnapshot(runner string, map_name string, map_points int) {
	var summary model.Summary
	// 取得剛才 UpdateGlobalSummary 更新後的最新資料
	if err := db.GetDB().Last(&summary).Error; err == nil {
		RecordGrowthSnapshot(summary.CurrentScore, summary.CompletedMaps, runner, map_name, map_points)
	}
}
