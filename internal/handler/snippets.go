package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

func RegisterSnippets(g *gin.RouterGroup) {
	g.GET("/snippets", func(c *gin.Context) {
		var items []model.Snippet
		database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"snippets": items})
	})
	g.GET("/snippets/templates", func(c *gin.Context) {
		var items []model.Snippet
		database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"templates": items})
	})
	g.POST("/snippets", func(c *gin.Context) {
		var item model.Snippet
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
		}
		database.DB.Create(&item)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "snippet": item})
	})
	g.PUT("/snippets/:id", func(c *gin.Context) {
		var item model.Snippet
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
		}
		c.ShouldBindJSON(&item)
		database.DB.Save(&item)
		c.JSON(http.StatusOK, gin.H{"ok": true, "snippet": item})
	})
	g.DELETE("/snippets/:id", func(c *gin.Context) {
		database.DB.Delete(&model.Snippet{}, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}