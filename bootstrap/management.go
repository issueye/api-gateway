package bootstrap

import (
	"fmt"

	"api-gateway/internal/api"
	"api-gateway/internal/global"
	"api-gateway/internal/services"
	"api-gateway/pkg/db"

	"github.com/gin-gonic/gin"
)

// ManagementApp结构体用于管理相关的初始化和配置
type ManagementApp struct {
	Router *gin.Engine
	API    *api.APIController
}

// NewManagementApp创建并初始化用于管理的应用实例
func NewManagementApp() *ManagementApp {
	return &ManagementApp{}
}

// Initialize初始化管理应用的各种组件
func (ma *ManagementApp) Initialize() {
	var err error
	global.DB, err = db.NewDB()
	if err != nil {
		fmt.Printf("Error initializing database: %v", err)
		return
	}

	apiService := services.NewAPIService()
	ma.API = api.NewAPIController(apiService)

	ma.Router = gin.Default()
}

// SetupRoutes设置管理应用的路由
func (ma *ManagementApp) SetupRoutes() {
	apiRoutes := ma.Router.Group("/api")
	{
		apiRoutes.POST("", ma.API.CreateAPI)
		apiRoutes.GET("", ma.API.GetAPIs)
		apiRoutes.GET("/:name", ma.API.GetAPI)
		apiRoutes.PUT("/:name", ma.API.UpdateAPI)
		apiRoutes.DELETE("/:name", ma.API.DeleteAPI)
	}
}

// Run启动管理应用
func (ma *ManagementApp) Run() {
	fmt.Println("Management API started on :8081")
	err := ma.Router.Run(":8081")
	if err != nil {
		fmt.Printf("Error starting management server: %v", err)
	}
}
