package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

func RegisterNamespaces(g *gin.RouterGroup) {
	g.GET("/namespaces", func(c *gin.Context) {
		var items []model.Namespace
		database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"namespaces": items})
	})
	g.POST("/namespaces", func(c *gin.Context) {
		var item model.Namespace
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
		}
		database.DB.Create(&item)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "namespace": item})
	})
	g.PUT("/namespaces/:id", func(c *gin.Context) {
		var item model.Namespace
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
		}
		c.ShouldBindJSON(&item)
		database.DB.Save(&item)
		c.JSON(http.StatusOK, gin.H{"ok": true, "namespace": item})
	})
	g.DELETE("/namespaces/:id", func(c *gin.Context) {
		database.DB.Delete(&model.Namespace{}, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}