package api

import (
	"context"
	"net/http"

	"api-gateway/internal/model"
	"api-gateway/internal/services"

	"github.com/gin-gonic/gin"
)

type DownstreamController struct {
	service services.DownstreamServiceImpl
}

func NewDownstreamController(service services.DownstreamServiceImpl) *DownstreamController {
	return &DownstreamController{
		service: service,
	}
}

// 创建API信息
func (ac *DownstreamController) Create(c *gin.Context) {
	var api model.Downstream
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
func (ac *DownstreamController) List(c *gin.Context) {
	result, err := ac.service.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

// 根据名称获取API信息
func (ac *DownstreamController) GetByName(c *gin.Context) {
	name := c.Param("name")
	api, err := ac.service.GetByName(context.Background(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, api)
}

// 更新API信息
func (ac *DownstreamController) Update(c *gin.Context) {
	name := c.Param("name")
	var data model.Downstream
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.service.UpdateByName(context.Background(), data, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// 删除API信息
func (ac *DownstreamController) Delete(c *gin.Context) {
	name := c.Param("name")
	err := ac.service.DeleteByName(context.Background(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
