package bootstrap

import (
	"fmt"

	"api-gateway/internal/api"
	"api-gateway/internal/services"

	"github.com/gin-gonic/gin"
)

// ManagementApp结构体用于管理相关的初始化和配置
type ManagementApp struct {
	Router     *gin.Engine
	API        *api.APIController
	DOWNStream *api.DownstreamController
}

// NewManagementApp创建并初始化用于管理的应用实例
func NewManagementApp() *ManagementApp {
	return &ManagementApp{}
}

// Initialize初始化管理应用的各种组件
func (ma *ManagementApp) Initialize() {
	apiService := services.NewAPIService()
	ma.API = api.NewAPIController(apiService)
	ma.Router = gin.Default()
}

// SetupRoutes设置管理应用的路由
func (ma *ManagementApp) SetupRoutes() {
	apiRoutes := ma.Router.Group("/api")
	{
		apiRoutes.POST("", ma.API.Create)
		apiRoutes.GET("", ma.API.List)
		apiRoutes.GET("/:name", ma.API.GetByName)
		apiRoutes.PUT("/:name", ma.API.Update)
		apiRoutes.DELETE("/:name", ma.API.Delete)
	}
	dsRoutes := ma.Router.Group("/downstream")
	{
		dsRoutes.POST("", ma.DOWNStream.Create)
		dsRoutes.GET("", ma.DOWNStream.List)
		dsRoutes.GET("/:name", ma.DOWNStream.GetByName)
		dsRoutes.PUT("/:name", ma.DOWNStream.Update)
		dsRoutes.DELETE("/:name", ma.DOWNStream.Delete)
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
