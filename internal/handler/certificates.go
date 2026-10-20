package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

func RegisterCertificates(g *gin.RouterGroup) {
	g.GET("/certificates", func(c *gin.Context) {
		var items []model.Certificate
		database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"certificates": items})
	})
	g.POST("/certificates", func(c *gin.Context) {
		var item model.Certificate
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
		}
		database.DB.Create(&item)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "certificate": item})
	})
	g.PUT("/certificates/:id", func(c *gin.Context) {
		var item model.Certificate
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
		}
		c.ShouldBindJSON(&item)
		database.DB.Save(&item)
		c.JSON(http.StatusOK, gin.H{"ok": true, "certificate": item})
	})
	g.DELETE("/certificates/:id", func(c *gin.Context) {
		database.DB.Delete(&model.Certificate{}, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}