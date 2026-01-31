package main

import (
	"log"

	"github.com/joho/godotenv"

	"DDNETONE/db"
	"DDNETONE/router"
	"DDNETONE/service"
)

func main() {
	// 1. 載入環境變數
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. 初始化資料庫
	db.Init()

	// 3. 啟動時更新一次全服總覽 (Optional, 視需求)
	service.UpdateGlobalSummary()

	// 4. 初始化路由並啟動 Server
	r := router.InitRouter()

	r.Run(":8080")
}
