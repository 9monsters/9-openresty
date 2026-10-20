package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

func RegisterStreams(g *gin.RouterGroup) {
	g.GET("/streams", func(c *gin.Context) {
		var items []model.Stream
		database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"streams": items})
	})
	g.POST("/streams", func(c *gin.Context) {
		var item model.Stream
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
		}
		database.DB.Create(&item)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "stream": item})
	})
	g.PUT("/streams/:id", func(c *gin.Context) {
		var item model.Stream
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
		}
		c.ShouldBindJSON(&item)
		database.DB.Save(&item)
		c.JSON(http.StatusOK, gin.H{"ok": true, "stream": item})
	})
	g.DELETE("/streams/:id", func(c *gin.Context) {
		database.DB.Delete(&model.Stream{}, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}