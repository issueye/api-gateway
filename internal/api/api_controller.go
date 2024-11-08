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
func (ac *APIController) Create(c *gin.Context) {
	var api model.APIInfo
	if err := c.ShouldBindJSON(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.service.Add(context.Background(), &api)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, api)
}

// 获取所有API信息
func (ac *APIController) List(c *gin.Context) {
	result, err := ac.service.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

// 根据名称获取API信息
func (ac *APIController) GetByName(c *gin.Context) {
	name := c.Param("name")
	api, err := ac.service.GetByName(context.Background(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, api)
}

// 更新API信息
func (ac *APIController) Update(c *gin.Context) {
	name := c.Param("name")
	var updatedAPI model.APIInfo
	if err := c.ShouldBindJSON(&updatedAPI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.service.UpdateByName(context.Background(), updatedAPI, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedAPI)
}

// 删除API信息
func (ac *APIController) Delete(c *gin.Context) {
	name := c.Param("name")
	err := ac.service.DeleteByName(context.Background(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
