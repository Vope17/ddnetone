package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// 2. 寫入資料庫
	// 注意：這裡會觸發我們之前寫的 BeforeSave Hook
	// 只要 Runner 有填且 Score > 0，Status 就會自動設為 2 (已完成)
	if err := db.Create(&newRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}

	c.JSON(http.StatusCreated, newRecord)

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
