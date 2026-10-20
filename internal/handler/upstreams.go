package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

func RegisterUpstreams(g *gin.RouterGroup) {
	g.GET("/upstreams", func(c *gin.Context) {
		var items []model.Upstream
		database.DB.Preload("Servers").Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"upstreams": items})
	})
	g.POST("/upstreams", func(c *gin.Context) {
		var item model.Upstream
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
		}
		database.DB.Create(&item)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "upstream": item})
	})
	g.PUT("/upstreams/:id", func(c *gin.Context) {
		var item model.Upstream
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
		}
		c.ShouldBindJSON(&item)
		database.DB.Save(&item)
		c.JSON(http.StatusOK, gin.H{"ok": true, "upstream": item})
	})
	g.DELETE("/upstreams/:id", func(c *gin.Context) {
		database.DB.Delete(&model.Upstream{}, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}