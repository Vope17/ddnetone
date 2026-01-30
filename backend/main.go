package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	ID         uint       `gorm:"primaryKey" json:"id"`
	Difficulty string     `json:"difficulty"`
	MapName    string     `json:"map_name"`
	Runner     string     `json:"runner"`
	Score      int        `json:"score"`
	Points     int        `json:"points"`
	Stars      int        `json:"stars"` // ★ 新增：星星數 (0-5)
	Note       string     `json:"note"`
	Status     int        `json:"status"` // 0:未完成, 1:進行中, 2:已完成
	FinishTime *time.Time `gorm:"column:finish_time" json:"finish_time"`
	HasDummy   bool       `gorm:"column:has_dummy" json:"has_dummy"`
}

// --- Models ---
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	User      string    `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// 自動判斷邏輯 & 自動填寫完成時間
func (m *MapRecord) BeforeSave(tx *gorm.DB) error {
	if m.Runner == "" || m.Runner == "-" || m.Runner == "nan" {
		m.Status = 0
	} else if m.Status == 1 {
		// ★ 修改：如果明確指定為 WIP (1)，則保持為 1，不自動變更為 2
	} else if m.Score > 0 {
		m.Status = 2
	} else {
		m.Status = 1
	}

	if m.Status == 2 && m.Score == 0 && m.Points > 0 {
		m.Score = m.Points
	}

	// ★ 新增：如果狀態是已完成，且還沒有記錄時間，則自動填入現在時間
	if m.Status == 2 && m.FinishTime == nil {
		now := time.Now()
		m.FinishTime = &now
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

	var err error
	err = godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbTimeZone := os.Getenv("DB_TIMEZONE")

	// 組合 DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbTimeZone)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&Summary{}, &Player{}, &MapRecord{}, &GrowthData{}, &Message{})

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

		api.GET("/messages", getMessages)
		api.POST("/messages", createMessage)
	}

	// seedData()

	// 新增：更新全服總覽
	updateGlobalSummary()

	r.Run(":8080")
}

// --- Handlers (保持不變) ---
func getSummary(c *gin.Context) {
	updateGlobalSummary()
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

	// 1. 計算 7 天前的時間點，並轉為 ISO 8601 字串格式
	// 這樣才能跟資料庫裡的字串欄位做比較
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)

	// 2. 加入 Where 條件：timestamp >= 7天前
	// 同時依照 id 或 timestamp 排序確保線條順序正確
	db.Where("timestamp >= ?", sevenDaysAgo).Order("id asc").Find(&growth)

	c.JSON(http.StatusOK, growth)
}

func createRecord(c *gin.Context) {
	var newRecord MapRecord

	// 1. 綁定前端傳來的 JSON
	if err := c.ShouldBindJSON(&newRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	// 2. 檢查地圖是否已存在 (避免重複新增)
	// 邏輯：先找找看有沒有同名且同難度且未完成的地圖
	var existingRecord MapRecord
	// 注意：這裡假設 Status != 2 代表未完成/進行中
	result := db.Where("map_name = ? AND difficulty = ? AND status != 2", newRecord.MapName, newRecord.Difficulty).First(&existingRecord)

	if result.Error == nil {

		// --- A. 更新現有紀錄 (Update) ---

		// 更新欄位
		existingRecord.Runner = newRecord.Runner

		existingRecord.Note = newRecord.Note
		existingRecord.FinishTime = &now
		existingRecord.Status = newRecord.Status
		existingRecord.HasDummy = newRecord.HasDummy
		existingRecord.Score = newRecord.Score

		// ★★★ 關鍵新增：更新跑者積分 ★★★
		if existingRecord.Status == 2 {
			updatePlayerStats(existingRecord.Runner, existingRecord.Score)
		} else {
			existingRecord.Score = 0
			existingRecord.HasDummy = false
		}

		// 儲存 (會觸發 BeforeSave 自動變更 Status 為 2)
		if err := db.Save(&existingRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}

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

		if newRecord.Score > 0 || newRecord.Status == 2 {
			newRecord.FinishTime = &now
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
				Role:              "PLAYER",
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

	recordGrowthSnapshot(int(totalScore), int(completedCount))
}

// ★ 新增的輔助函式：寫入成長紀錄
func recordGrowthSnapshot(score int, maps int) {
	// 1. 找出第一筆完成紀錄的時間 (計算 Hours 用)
	var firstRecord MapRecord
	var startTime time.Time

	err := db.Where("status = 2 AND finish_time IS NOT NULL").Order("finish_time asc").First(&firstRecord).Error
	if err == nil && firstRecord.FinishTime != nil {
		startTime = *firstRecord.FinishTime
	} else {
		startTime = time.Now()
	}

	hoursSinceStart := time.Since(startTime).Hours()
	if hoursSinceStart < 0 {
		hoursSinceStart = 0
	}

	// ★★★ 新增：檢查上一筆資料，避免重複記錄 ★★★
	var lastGrowth GrowthData
	// 找出最新的一筆成長紀錄
	if err := db.Order("id desc").First(&lastGrowth).Error; err == nil {
		// 如果「分數」和「地圖數」都沒變，就直接退出，不存入資料庫
		if lastGrowth.Points == score && lastGrowth.Maps == maps {
			return
		}
	}

	// 3. 寫入 GrowthData (只有數值改變時才會執行到這)
	newGrowth := GrowthData{
		Hours:     hoursSinceStart,
		Points:    score,
		Maps:      maps,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	db.Create(&newGrowth)
}

// --- Handlers ---
func getMessages(c *gin.Context) {
	var messages []Message
	// 依時間降冪排列，讓最新的留言在上面
	db.Order("created_at desc").Find(&messages)
	c.JSON(http.StatusOK, messages)
}

func createMessage(c *gin.Context) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg.CreatedAt = time.Now()
	if err := db.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post message"})
		return
	}
	c.JSON(http.StatusCreated, msg)
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
