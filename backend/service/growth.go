package service

import (
	"net/http"
	"time"

	"DDNETONE/db"
	"DDNETONE/model"
	"github.com/gin-gonic/gin"
)

func GetGrowth(c *gin.Context) {
	var growth []model.GrowthData
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	db.GetDB().Where("timestamp >= ?", sevenDaysAgo).Order("id asc").Find(&growth)
	c.JSON(http.StatusOK, growth)
}

type DailyActivity struct {
	Date  string `json:"date"`
	Maps  int    `json:"maps"`
	Score int    `json:"score"`
}

// GetDailyActivity 回傳每日完成地圖數與分數增量（過去一年）
func GetDailyActivity(c *gin.Context) {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	// 取得完成記錄，按日期彙整
	type row struct {
		Date  string
		Maps  int
		Score int
	}

	var rows []row
	db.GetDB().Model(&model.MapRecord{}).
		Select("DATE(finish_time AT TIME ZONE 'Asia/Taipei') AS date, COUNT(*) AS maps, SUM(score) AS score").
		Where("status = 2 AND finish_time >= ?", oneYearAgo).
		Group("DATE(finish_time AT TIME ZONE 'Asia/Taipei')").
		Order("date asc").
		Scan(&rows)

	result := make([]DailyActivity, len(rows))
	for i, r := range rows {
		result[i] = DailyActivity{Date: r.Date, Maps: r.Maps, Score: r.Score}
	}

	c.JSON(http.StatusOK, result)
}

type MilestoneResult struct {
	Target    int    `json:"target"`
	Timestamp string `json:"timestamp"`
	Maps      int    `json:"maps"`
}

func GetMilestones(c *gin.Context) {
	const TARGET_MAPS = 2403

	// 單次查詢，按 maps asc 排序後線性掃描，避免 N 次查詢
	var records []model.GrowthData
	db.GetDB().Order("maps asc, id asc").Find(&records)

	results := []MilestoneResult{}
	target := 100
	for _, r := range records {
		for target <= TARGET_MAPS && r.Maps >= target {
			results = append(results, MilestoneResult{
				Target:    target,
				Timestamp: r.Timestamp,
				Maps:      r.Maps,
			})
			target += 100
		}
		if target > TARGET_MAPS {
			break
		}
	}

	c.JSON(http.StatusOK, results)
}

func GetScoreMilestones(c *gin.Context) {
	const SCORE_STEP = 1000

	// 單次查詢，按 points asc 排序後線性掃描，避免 N 次查詢
	var records []model.GrowthData
	db.GetDB().Order("points asc, id asc").Find(&records)

	results := []MilestoneResult{}
	target := SCORE_STEP
	for _, r := range records {
		for r.Points >= target {
			results = append(results, MilestoneResult{
				Target:    target,
				Timestamp: r.Timestamp,
				Maps:      r.Points,
			})
			target += SCORE_STEP
		}
	}

	c.JSON(http.StatusOK, results)
}

func RecordGrowthSnapshot(score int, maps int, runner string, map_name string, map_points int) {
	database := db.GetDB()
	var firstRecord model.MapRecord
	var startTime time.Time

	err := database.Where("status = 2 AND finish_time IS NOT NULL").Order("finish_time asc").First(&firstRecord).Error
	if err == nil && firstRecord.FinishTime != nil {
		startTime = *firstRecord.FinishTime
	} else {
		startTime = time.Now()
	}

	hoursSinceStart := time.Since(startTime).Hours()
	if hoursSinceStart < 0 {
		hoursSinceStart = 0
	}

	var lastGrowth model.GrowthData
	if err := database.Order("id desc").First(&lastGrowth).Error; err == nil {
		if lastGrowth.Points == score && lastGrowth.Maps == maps {
			return
		}
	}

	newGrowth := model.GrowthData{
		Hours:     hoursSinceStart,
		Points:    score,
		Maps:      maps,
		Runner:    runner,
		MapName:   map_name,
		MapPoints: map_points,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	database.Create(&newGrowth)
}
