package bootstrap

import (
	"api-gateway/internal/global"
	"api-gateway/internal/model"
	"api-gateway/pkg/db"
	"api-gateway/pkg/logger"
	"fmt"
	"path/filepath"
)

func Run() {
	// 初始化运行路径
	InitRuntime()
	// 初始化日志
	InitLogger()
	// 初始化数据库
	InitDB()
	// 启动网关服务
	RunGetWay()
	// 启动管理服务
	RunManagement()
	select {}
}

func RunGetWay() {
	go func() {
		gatewayApp := NewGatewayApp()
		defer gatewayApp.Close()

		gatewayApp.Initialize()
		gatewayApp.SetupRoutes()
		gatewayApp.Run()
	}()
}

func RunManagement() {
	go func() {
		managementApp := NewManagementApp()
		managementApp.Initialize()
		managementApp.SetupRoutes()
		managementApp.Run()
	}()
}

func InitLogger() {
	path := filepath.Join(LOG_PATH, "server.log")
	global.Logger = logger.InitLogger(path)
}

func InitDB() {
	path := filepath.Join(DB_PATH, "data.db")
	var err error
	global.DB, err = db.NewDB(path)
	if err != nil {
		panic(fmt.Sprintf("初始化[sqlite]数据库失败: %v", err))
	}

	global.DB.AutoMigrate(
		&model.APIInfo{},
		&model.Downstream{},
		&model.TrafficStats{},
	)
}
