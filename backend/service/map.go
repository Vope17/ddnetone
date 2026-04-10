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

		c.JSON(http.StatusOK, existingRecord)

	} else {
		// --- B. Create ---
		if err := database.Create(&newRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
			return
		}

		UpdateGlobalSummary()
		triggerSnapshot(newRecord.Runner, newRecord.MapName, newRecord.Score)

		c.JSON(http.StatusCreated, newRecord)
	}
}

func GetLoadOptions(c *gin.Context) {
	var maps []model.MapRecord
	db.GetDB().Where("status = 2").Order("map_name asc").Find(&maps)
	c.JSON(http.StatusOK, maps)
}

type LoadRecordRequest struct {
	MapName    string `json:"map_name" binding:"required"`
	Difficulty string `json:"difficulty" binding:"required"`
}

func LoadRecord(c *gin.Context) {
	var req LoadRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database := db.GetDB()

	// 確認目前總分已達 10000
	var completedScore int64
	var loadedScore int64
	database.Model(&model.MapRecord{}).Where("status = 2").Select("COALESCE(SUM(points), 0)").Scan(&completedScore)
	database.Model(&model.MapRecord{}).Where("status = 3").Select("COALESCE(SUM(points), 0)").Scan(&loadedScore)
	if completedScore+loadedScore < 10000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "total score has not reached 10000 yet"})
		return
	}

	var record model.MapRecord
	if err := database.Where("map_name = ? AND difficulty = ? AND status = 2", req.MapName, req.Difficulty).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "completed record not found"})
		return
	}

	record.Status = 3
	if err := database.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load record"})
		return
	}

	UpdateGlobalSummary()
	triggerSnapshot(record.Runner, record.MapName, -record.Points)

	c.JSON(http.StatusOK, record)
}


func triggerSnapshot(runner string, map_name string, map_points int) {
	var summary model.Summary
	// 取得剛才 UpdateGlobalSummary 更新後的最新資料
	if err := db.GetDB().Last(&summary).Error; err == nil {
		RecordGrowthSnapshot(summary.CurrentScore, summary.CompletedMaps, runner, map_name, map_points)
	}
}
