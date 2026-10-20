package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

func RegisterAccessLists(g *gin.RouterGroup) {
	g.GET("/access-lists", func(c *gin.Context) {
		var items []model.AccessList
		database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"lists": items})
	})
	g.POST("/access-lists", func(c *gin.Context) {
		var item model.AccessList
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
		}
		database.DB.Create(&item)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "list": item})
	})
	g.PUT("/access-lists/:id", func(c *gin.Context) {
		var item model.AccessList
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
		}
		c.ShouldBindJSON(&item)
		database.DB.Save(&item)
		c.JSON(http.StatusOK, gin.H{"ok": true, "list": item})
	})
	g.DELETE("/access-lists/:id", func(c *gin.Context) {
		database.DB.Delete(&model.AccessList{}, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}