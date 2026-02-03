package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"DDNETONE/service" // 引入 service
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"http://localhost:5173"},

		AllowAllOrigins: true,

		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	api := r.Group("/api")
	{

		api.GET("/summary", service.GetSummary)

		api.GET("/leaderboard", service.GetLeaderboard)

		api.GET("/maps", service.GetMaps)
		api.POST("/records", service.CreateRecord)
		api.GET("/map-options", service.GetMapOptions)

		api.GET("/player-options", service.GetPlayerOptions)

		api.GET("/growth", service.GetGrowth)

		api.GET("/messages", service.GetMessages)
		api.POST("/messages", service.CreateMessage)

	}

	return r
}
