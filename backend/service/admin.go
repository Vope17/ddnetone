package service

import (
	"net/http"
	"os"

	"DDNETONE/db"
	"DDNETONE/model"
	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 驗證 X-Admin-Key header
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Admin-Key")
		if key == "" || key != os.Getenv("ADMIN_KEY") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

// GetAdminRecords 取得所有已完成及已加載記錄供管理
func GetAdminRecords(c *gin.Context) {
	var records []model.MapRecord
	db.GetDB().Where("status IN (2, 3)").Order("finish_time desc").Find(&records)
	c.JSON(http.StatusOK, records)
}

type EditRecordRequest struct {
	Note   *string `json:"note"`
	Runner *string `json:"runner"`
}

// EditRecord 修改 note 或 runner
func EditRecord(c *gin.Context) {
	id := c.Param("id")
	var record model.MapRecord
	if err := db.GetDB().First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	var req EditRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database := db.GetDB()

	if req.Note != nil {
		record.Note = *req.Note
	}

	if req.Runner != nil {
		record.Runner = *req.Runner
	}

	if err := database.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update record"})
		return
	}

	BroadcastUpdate()
	c.JSON(http.StatusOK, record)
}

// UndoRecord 將已完成的記錄(status=2)還原為未完成(status=0)
func UndoRecord(c *gin.Context) {
	id := c.Param("id")
	var record model.MapRecord
	if err := db.GetDB().First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	if record.Status != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record is not completed"})
		return
	}

	database := db.GetDB()

	// 刪除對應的 growth_data 快照
	database.Where("map_name = ? AND runner = ?", record.MapName, record.Runner).
		Delete(&model.GrowthData{})

	// 重置記錄
	record.Status = 0
	record.Runner = ""
	record.Score = 0
	record.FinishTime = nil
	record.HasDummy = false

	if err := database.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to undo record"})
		return
	}

	UpdateGlobalSummary()
	BroadcastUpdate()
	c.JSON(http.StatusOK, record)
}

type CreateAdminMapRequest struct {
	MapName    string `json:"map_name" binding:"required"`
	Difficulty string `json:"difficulty" binding:"required"`
	Points     int    `json:"points"`
	Stars      int    `json:"stars"`
}

// UnloadRecord 將已加載的記錄(status=3)還原為已完成(status=2)
func UnloadRecord(c *gin.Context) {
	id := c.Param("id")
	var record model.MapRecord
	if err := db.GetDB().First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	if record.Status != 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record is not loaded"})
		return
	}

	database := db.GetDB()

	// 刪除對應的 growth_data 快照
	database.Where("map_name = ? AND runner = ?", record.MapName, record.Runner).
		Delete(&model.GrowthData{})

	record.Status = 2

	if err := database.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unload record"})
		return
	}

	UpdateGlobalSummary()
	BroadcastUpdate()
	c.JSON(http.StatusOK, record)
}

func CreateAdminMap(c *gin.Context) {
	var req CreateAdminMapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database := db.GetDB()
	var existing model.MapRecord
	if err := database.Where("map_name = ? AND difficulty = ?", req.MapName, req.Difficulty).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "map already exists"})
		return
	}

	record := model.MapRecord{
		MapName:    req.MapName,
		Difficulty: req.Difficulty,
		Points:     req.Points,
		Stars:      req.Stars,
		Status:     0,
	}

	if err := database.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create map"})
		return
	}

	c.JSON(http.StatusCreated, record)
}
