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

		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "X-Admin-Key"},
	}))

	api := r.Group("/api")
	{

		api.GET("/summary", service.GetSummary)

		api.GET("/leaderboard", service.GetLeaderboard)

		api.GET("/maps", service.GetMaps)
		api.POST("/records", service.CreateRecord)
		api.GET("/map-options", service.GetMapOptions)
		api.GET("/load-options", service.GetLoadOptions)
		api.POST("/records/load", service.LoadRecord)

		api.GET("/player-options", service.GetPlayerOptions)

		api.GET("/growth", service.GetGrowth)
		api.GET("/milestones", service.GetMilestones)
		api.GET("/score-milestones", service.GetScoreMilestones)
		api.GET("/daily-activity", service.GetDailyActivity)

		admin := api.Group("/admin")
		admin.Use(service.AdminAuthMiddleware())
		{
			admin.GET("/records", service.GetAdminRecords)
			admin.PUT("/records/:id", service.EditRecord)
			admin.PUT("/records/:id/undo", service.UndoRecord)
			admin.PUT("/records/:id/unload", service.UnloadRecord)
			admin.POST("/maps", service.CreateAdminMap)
		}

		api.GET("/messages", service.GetMessages)
		api.POST("/messages", service.CreateMessage)

	}

	return r
}
