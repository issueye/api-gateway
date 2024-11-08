package api

import (
	"context"
	"net/http"

	"api-gateway/internal/model"
	"api-gateway/internal/services"

	"github.com/gin-gonic/gin"
)

type APIController struct {
	service services.APIServiceImpl
}

func NewAPIController(service services.APIServiceImpl) *APIController {
	return &APIController{
		service: service,
	}
}

// 创建API信息
func (ac *APIController) CreateAPI(c *gin.Context) {
	var api model.APIInfo
	if err := c.ShouldBindJSON(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.service.AddAPIInfo(context.Background(), &api)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, api)
}

// 获取所有API信息
func (ac *APIController) GetAPIs(c *gin.Context) {
	result, err := ac.service.GetAllAPIInfo(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

// 根据名称获取API信息
func (ac *APIController) GetAPI(c *gin.Context) {
	name := c.Param("name")
	api, err := ac.service.GetAPIInfos(context.Background(), map[string]any{"name": name})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, api)
}

// 更新API信息
func (ac *APIController) UpdateAPI(c *gin.Context) {
	name := c.Param("name")
	var updatedAPI model.APIInfo
	if err := c.ShouldBindJSON(&updatedAPI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.service.UpdateAPIInfoByName(context.Background(), updatedAPI, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedAPI)
}

// 删除API信息
func (ac *APIController) DeleteAPI(c *gin.Context) {
	name := c.Param("name")
	err := ac.service.DeleteAPIInfo(context.Background(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
