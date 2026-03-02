package service

import (
	"net/http"
	"os"

	"DDNETONE/db"
	"DDNETONE/model"
	"DDNETONE/utils"
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

// GetAdminRecords 取得所有已完成記錄供管理
func GetAdminRecords(c *gin.Context) {
	var records []model.MapRecord
	db.GetDB().Where("status = 2").Order("finish_time desc").Find(&records)
	c.JSON(http.StatusOK, records)
}

type EditRecordRequest struct {
	Note   *string `json:"note"`
	Runner *string `json:"runner"`
}

// EditRecord 修改 note 或 runner（runner 改變時調整玩家積分）
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

	if req.Runner != nil && *req.Runner != record.Runner {
		// 從舊 runner 扣分
		reversePlayerStats(record.Runner, record.Score)
		// 給新 runner 加分
		UpdatePlayerStats(*req.Runner, record.Score)
		record.Runner = *req.Runner
	}

	if err := database.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update record"})
		return
	}

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

	// 反轉玩家積分
	reversePlayerStats(record.Runner, record.Score)

	// 刪除對應的 growth_data 快照（由此地圖完成所觸發的）
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
	c.JSON(http.StatusOK, record)
}

// GetAdminPlayers 取得所有玩家資訊供管理
func GetAdminPlayers(c *gin.Context) {
	var players []model.Player
	db.GetDB().Order("score_contribution desc").Find(&players)
	c.JSON(http.StatusOK, players)
}

type EditPlayerRequest struct {
	Name              *string  `json:"name"`
	Role              *string  `json:"role"`
	ScoreContribution *float64 `json:"score_contrib"`
	MapCount          *int     `json:"map_count"`
}

// EditPlayer 修改玩家資訊
func EditPlayer(c *gin.Context) {
	id := c.Param("id")
	var player model.Player
	if err := db.GetDB().First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	var req EditPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		player.Name = *req.Name
	}
	if req.Role != nil {
		player.Role = *req.Role
	}
	if req.ScoreContribution != nil {
		player.ScoreContribution = *req.ScoreContribution
	}
	if req.MapCount != nil {
		player.MapCount = *req.MapCount
	}

	if err := db.GetDB().Save(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update player"})
		return
	}

	c.JSON(http.StatusOK, player)
}

// DeletePlayer 刪除玩家
func DeletePlayer(c *gin.Context) {
	id := c.Param("id")
	var player model.Player
	if err := db.GetDB().First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}
	if err := db.GetDB().Delete(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete player"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "player deleted"})
}

// reversePlayerStats 從玩家積分中扣除指定分數
func reversePlayerStats(runnerNamesRaw string, score int) {
	if runnerNamesRaw == "" || score <= 0 {
		return
	}

	names := utils.ParseRunnerNames(runnerNamesRaw)
	gormDB := db.GetDB()

	for _, name := range names {
		var player model.Player
		if err := gormDB.Where("name = ?", name).First(&player).Error; err != nil {
			continue
		}
		player.ScoreContribution -= float64(score)
		if player.ScoreContribution < 0 {
			player.ScoreContribution = 0
		}
		player.MapCount -= 1
		if player.MapCount < 0 {
			player.MapCount = 0
		}
		gormDB.Save(&player)
	}
}
