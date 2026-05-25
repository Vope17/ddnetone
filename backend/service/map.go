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
	q := db.GetDB()
	if difficulty != "" && difficulty != "ALL" {
		q = q.Where("difficulty = ? AND status != 2", difficulty)
	} else {
		q = q.Where("status != 2")
	}
	q.Order("map_name asc").Find(&maps)
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

		if err := database.Save(&existingRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}

		UpdateGlobalSummary()
		// 呼叫 Growth Snapshot (需先取得最新 Summary 數據)
		triggerSnapshot(existingRecord.Runner, existingRecord.MapName, existingRecord.Score)

		BroadcastUpdate()
		c.JSON(http.StatusOK, existingRecord)

	} else {
		// --- B. Create ---
		if err := database.Create(&newRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
			return
		}

		UpdateGlobalSummary()
		triggerSnapshot(newRecord.Runner, newRecord.MapName, newRecord.Score)

		BroadcastUpdate()
		c.JSON(http.StatusCreated, newRecord)
	}
}

func triggerSnapshot(runner string, map_name string, map_points int) {
	var summary model.Summary
	// 取得剛才 UpdateGlobalSummary 更新後的最新資料
	if err := db.GetDB().Last(&summary).Error; err == nil {
		RecordGrowthSnapshot(summary.CurrentScore, summary.CompletedMaps, runner, map_name, map_points)
	}
}
