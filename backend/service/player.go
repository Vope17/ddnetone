package service

import (
	"net/http"

	"DDNETONE/db"
	"DDNETONE/model"
	"github.com/gin-gonic/gin"

	"DDNETONE/utils"
)

func GetLeaderboard(c *gin.Context) {
	var players []model.Player
	db.GetDB().Order("score_contribution desc").Find(&players)
	c.JSON(http.StatusOK, players)
}

// UpdatePlayerStats 更新跑者積分 (內部呼叫用)
func UpdatePlayerStats(runnerNamesRaw string, score int) {
	if runnerNamesRaw == "" || score <= 0 {
		return
	}

	names := utils.ParseRunnerNames(runnerNamesRaw)
	database := db.GetDB()

	for _, runnerName := range names {
		var player model.Player
		result := database.Where("name = ?", runnerName).First(&player)

		if result.Error == nil {
			player.ScoreContribution += float64(score)
			player.MapCount += 1
			database.Save(&player)
		} else {
			newPlayer := model.Player{

				Name: runnerName,

				Role:              "PLAYER",
				ScoreContribution: float64(score),
				MapCount:          1,
				ContributionRate:  0,
			}
			database.Create(&newPlayer)
		}
	}
}
