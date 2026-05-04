package service

import (
	"net/http"
	"sort"

	"DDNETONE/db"
	"DDNETONE/model"
	"DDNETONE/utils"
	"github.com/gin-gonic/gin"
)

// buildLeaderboard computes the leaderboard from DB (shared by API and SSE).
func buildLeaderboard() []model.PlayerStats {
	var records []model.MapRecord
	if err := db.GetDB().Where("status = 2 AND score > 0").Find(&records).Error; err != nil {
		return []model.PlayerStats{}
	}

	scoreMap := make(map[string]float64)
	countMap := make(map[string]int)
	for _, r := range records {
		for _, name := range utils.ParseRunnerNames(r.Runner) {
			scoreMap[name] += float64(r.Score)
			countMap[name]++
		}
	}

	var players []model.Player
	db.GetDB().Find(&players)
	roleMap := make(map[string]string)
	idMap := make(map[string]uint)
	for _, p := range players {
		roleMap[p.Name] = p.Role
		idMap[p.Name] = p.ID
	}

	var result []model.PlayerStats
	for name, score := range scoreMap {
		result = append(result, model.PlayerStats{
			ID:                idMap[name],
			Name:              name,
			Role:              roleMap[name],
			ScoreContribution: score,
			MapCount:          countMap[name],
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ScoreContribution > result[j].ScoreContribution
	})

	return result
}

func GetLeaderboard(c *gin.Context) {
	c.JSON(http.StatusOK, buildLeaderboard())
}

func GetPlayerOptions(c *gin.Context) {
	var names []string

	err := db.GetDB().
		Model(&model.Player{}).
		Distinct("name").
		Order("name ASC").
		Pluck("name", &names).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取玩家列表"})
		return
	}

	c.JSON(http.StatusOK, names)
}
