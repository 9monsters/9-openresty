package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

func RegisterAcme(g *gin.RouterGroup) {
	g.GET("/acme-users", func(c *gin.Context) {
		var items []model.AcmeUser
		database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"users": items})
	})
	g.POST("/acme-users", func(c *gin.Context) {
		var item model.AcmeUser
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
		}
		database.DB.Create(&item)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "user": item})
	})
	g.PUT("/acme-users/:id", func(c *gin.Context) {
		var item model.AcmeUser
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
		}
		c.ShouldBindJSON(&item)
		database.DB.Save(&item)
		c.JSON(http.StatusOK, gin.H{"ok": true, "user": item})
	})
	g.DELETE("/acme-users/:id", func(c *gin.Context) {
		database.DB.Delete(&model.AcmeUser{}, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}