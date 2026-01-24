package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

// --- Models ---

type Summary struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CurrentScore  int       `json:"current_score"`
	TargetScore   int       `json:"target_score"`
	CompletedMaps int       `json:"completed_maps"`
	LastUpdate    time.Time `json:"last_update"`
}

type Player struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	Name              string  `json:"name"`
	Role              string  `json:"role"`
	ScoreContribution float64 `json:"score_contrib"`
	MapCount          int     `json:"map_count"`
	ContributionRate  float64 `json:"contribution_rate"`
}

type MapRecord struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Difficulty string `json:"difficulty"`
	MapName    string `json:"map_name"`
	Runner     string `json:"runner"`
	Score      int    `json:"score"`
	Points     int    `json:"points"`
	Note       string `json:"note"`
	Status     int    `json:"status"` // 0:未完成, 1:進行中, 2:已完成
}

// ★★★ 這裡加入了自動判斷邏輯 ★★★
func (m *MapRecord) BeforeSave(tx *gorm.DB) error {
	if m.Runner == "" || m.Runner == "-" || m.Runner == "nan" {
		m.Status = 0 // Incomplete
	} else if m.Score > 0 {
		m.Status = 2 // Completed
	} else {
		m.Status = 1 // In Progress
	}

	if m.Status == 2 && m.Score == 0 && m.Points > 0 {
		m.Score = m.Points
	}

	return nil
}

type GrowthData struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Hours     float64 `json:"hours"`
	Points    int     `json:"points"`
	Maps      int     `json:"maps"`
	Timestamp string  `json:"timestamp"`
}

var db *gorm.DB

func main() {

	// 請確認密碼與 DB 名稱
	dsn := "host=localhost user=postgres password=123456 dbname=ddnetone port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&Summary{}, &Player{}, &MapRecord{}, &GrowthData{})

	r := gin.Default()
	// r.Static("/assets", "./dist/assets")
	// r.LoadHTMLGlob("dist/*.html")

	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", nil)
	// })

	r.Use(cors.New(cors.Config{
		// AllowAllOrigins: true,
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	api := r.Group("/api")
	{
		api.GET("/summary", getSummary)

		api.GET("/leaderboard", getLeaderboard)
		api.GET("/maps", getMaps)
		api.GET("/growth", getGrowth)
		api.POST("/records", createRecord)
		api.GET("/map-options", getMapOptions)
	}

	seedData()

	// ★★★ 新增：更新全服總覽 ★★★
	updateGlobalSummary()

	r.Run(":8080")
}

// --- Handlers (保持不變) ---
func getSummary(c *gin.Context) {
	var summary Summary
	db.Last(&summary)
	c.JSON(http.StatusOK, summary)
}

func getLeaderboard(c *gin.Context) {
	var players []Player

	// 直接撈取 Player 表，並按分數排序
	db.Order("score_contribution desc").Find(&players)
	c.JSON(http.StatusOK, players)
}

func getMaps(c *gin.Context) {

	difficulty := c.Query("difficulty")
	var maps []MapRecord
	query := db.Model(&MapRecord{})
	if difficulty != "" {

		query = query.Where("difficulty = ?", difficulty)
	}
	// 按照分數排序，讓完成的在上面
	query.Order("score desc").Find(&maps)
	c.JSON(http.StatusOK, maps)
}

func getGrowth(c *gin.Context) {
	var growth []GrowthData

	db.Order("hours asc").Find(&growth)

	c.JSON(http.StatusOK, growth)
}
func createRecord(c *gin.Context) {
	var newRecord MapRecord

	// 1. 綁定前端傳來的 JSON
	if err := c.ShouldBindJSON(&newRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 檢查地圖是否已存在 (避免重複新增)
	// 邏輯：先找找看有沒有同名且同難度且未完成的地圖
	var existingRecord MapRecord
	// 注意：這裡假設 Status != 2 代表未完成/進行中
	result := db.Where("map_name = ? AND difficulty = ? AND status != 2", newRecord.MapName, newRecord.Difficulty).First(&existingRecord)

	if result.Error == nil {

		// --- A. 更新現有紀錄 (Update) ---

		// 更新欄位
		existingRecord.Runner = newRecord.Runner

		existingRecord.Score = newRecord.Score // 這是玩家獲得的分數
		existingRecord.Note = newRecord.Note

		// 儲存 (會觸發 BeforeSave 自動變更 Status 為 2)
		if err := db.Save(&existingRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}

		// ★★★ 關鍵新增：更新跑者積分 ★★★
		updatePlayerStats(existingRecord.Runner, existingRecord.Score)

		// ★★★ 新增：更新全服總覽 ★★★
		updateGlobalSummary()

		c.JSON(http.StatusOK, existingRecord)

	} else {
		// --- B. 新增全新紀錄 (Create) ---

		if err := db.Create(&newRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
			return
		}

		// ★★★ 關鍵新增：更新跑者積分 ★★★
		// 只有當真的完成了 (Status == 2) 才加分
		if newRecord.Status == 2 {
			updatePlayerStats(newRecord.Runner, newRecord.Score)
		}

		// ★★★ 新增：更新全服總覽 ★★★
		updateGlobalSummary()

		c.JSON(http.StatusCreated, newRecord)
	}
}

// --- 輔助函式：更新跑者積分 (支援多人，用逗號或 & 分隔) ---
func updatePlayerStats(runnerNamesRaw string, score int) {
	if runnerNamesRaw == "" || score <= 0 {
		return
	}

	// 1. 處理分隔符號：支援 "A, B" 或 "A & B" 或 "A,B"
	// 先把 '&' 替換成 ','，再統一用 ',' 切割
	normalized := strings.ReplaceAll(runnerNamesRaw, "&", ",")
	names := strings.Split(normalized, ",")

	for _, name := range names {
		// 去除前後空白 (例如 " PlayerB" 變成 "PlayerB")
		runnerName := strings.TrimSpace(name)

		if runnerName == "" {
			continue
		}

		// --- 以下是原本的單人更新邏輯，現在包在迴圈裡對每個人執行 ---
		var player Player
		result := db.Where("name = ?", runnerName).First(&player)

		if result.Error == nil {
			// (A) 跑者已存在 -> 更新
			player.ScoreContribution += float64(score)
			player.MapCount += 1
			db.Save(&player)
		} else {
			// (B) 跑者不存在 -> 建立
			newPlayer := Player{
				Name:              runnerName,
				Role:              "Agent",
				ScoreContribution: float64(score),
				MapCount:          1,
				ContributionRate:  0,
			}
			db.Create(&newPlayer)
		}
	}
}

// 取得下拉選單用的地圖列表
func getMapOptions(c *gin.Context) {
	difficulty := c.Query("difficulty")

	var maps []MapRecord

	// 直接根據 Status 篩選
	// 找出該難度下，Status 不是 2 (已完成) 的地圖
	db.Where("difficulty = ? AND status != 2", difficulty).
		Order("map_name asc").
		Find(&maps)

	c.JSON(http.StatusOK, maps)
}

// --- 輔助函式：更新全服總覽數據 ---
func updateGlobalSummary() {
	var totalScore int64
	var completedCount int64

	// 1. 計算所有「已完成 (Status=2)」地圖的分數總和 (使用 Points 欄位)
	// COALESCE 確保如果沒有資料時回傳 0 而不是 null
	db.Model(&MapRecord{}).Where("status = 2").Select("COALESCE(SUM(points), 0)").Scan(&totalScore)

	// 2. 計算已完成的地圖數量
	db.Model(&MapRecord{}).Where("status = 2").Count(&completedCount)

	// 3. 更新 Summary 表
	var summary Summary
	// 嘗試抓取最後一筆 Summary (通常只有一筆)
	if err := db.Last(&summary).Error; err != nil {
		// 如果還沒有 Summary，就建立一筆新的
		summary = Summary{
			TargetScore: 32450, // 您的目標分數
		}
		db.Create(&summary)
	}

	// 更新數值
	summary.CurrentScore = int(totalScore)
	summary.CompletedMaps = int(completedCount)
	summary.LastUpdate = time.Now()

	// 寫回資料庫
	db.Save(&summary)
}

// --- Seed Data ---

func seedData() {
	var count int64
	db.Model(&MapRecord{}).Count(&count)
	if count == 0 {
		log.Println("Seeding data...")

		// 1. 測試已完成 (Status 會自動變 2)

		db.Create(&MapRecord{Difficulty: "Insane", MapName: "Ravillion", Runner: "MEE6", Score: 38, Note: "First clear"})

		// 2. 測試進行中 (有 Runner 但 Score 為 0 -> Status 會自動變 1)
		db.Create(&MapRecord{Difficulty: "Insane", MapName: "ImpossibleMap", Runner: "Hardy", Score: 0, Note: "Trying..."})

		// 3. 測試未完成 (無 Runner -> Status 會自動變 0)
		db.Create(&MapRecord{Difficulty: "Novice", MapName: "JustFly", Runner: "", Score: 0, Note: ""})

		// 其他資料...
		db.Create(&Summary{CurrentScore: 4202, TargetScore: 32450, CompletedMaps: 401, LastUpdate: time.Now()})
		db.Create(&Player{Name: "YuanYuan", Role: "主持人", ScoreContribution: 1037, MapCount: 47, ContributionRate: 0.24})
		db.Create(&GrowthData{Hours: 6, Points: 1037, Maps: 89})
	}
}
