package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

func RegisterDNS(g *gin.RouterGroup) {
	g.GET("/dns-credentials", func(c *gin.Context) {
		var items []model.DNSCredential; database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"credentials": items})
	})
	g.POST("/dns-credentials", func(c *gin.Context) {
		var item model.DNSCredential
		if err := c.ShouldBindJSON(&item); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return }
		database.DB.Create(&item); c.JSON(http.StatusCreated, gin.H{"ok": true, "credential": item})
	})
	g.PUT("/dns-credentials/:id", func(c *gin.Context) {
		var item model.DNSCredential
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return }
		c.ShouldBindJSON(&item); database.DB.Save(&item); c.JSON(http.StatusOK, gin.H{"ok": true, "credential": item})
	})
	g.DELETE("/dns-credentials/:id", func(c *gin.Context) { database.DB.Delete(&model.DNSCredential{}, c.Param("id")); c.JSON(http.StatusOK, gin.H{"ok": true}) })

	g.GET("/dns-domains", func(c *gin.Context) {
		var items []model.DNSDomain; database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"domains": items})
	})
	g.POST("/dns-domains", func(c *gin.Context) {
		var item model.DNSDomain
		if err := c.ShouldBindJSON(&item); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return }
		database.DB.Create(&item); c.JSON(http.StatusCreated, gin.H{"ok": true, "domain": item})
	})
	g.PUT("/dns-domains/:id", func(c *gin.Context) {
		var item model.DNSDomain
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return }
		c.ShouldBindJSON(&item); database.DB.Save(&item); c.JSON(http.StatusOK, gin.H{"ok": true, "domain": item})
	})
	g.DELETE("/dns-domains/:id", func(c *gin.Context) { database.DB.Delete(&model.DNSDomain{}, c.Param("id")); c.JSON(http.StatusOK, gin.H{"ok": true}) })

	g.GET("/dns-groups", func(c *gin.Context) {
		var items []model.DNSGroup; database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"groups": items})
	})
	g.POST("/dns-groups", func(c *gin.Context) {
		var item model.DNSGroup
		if err := c.ShouldBindJSON(&item); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return }
		database.DB.Create(&item); c.JSON(http.StatusCreated, gin.H{"ok": true, "group": item})
	})
	g.PUT("/dns-groups/:id", func(c *gin.Context) {
		var item model.DNSGroup
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return }
		c.ShouldBindJSON(&item); database.DB.Save(&item); c.JSON(http.StatusOK, gin.H{"ok": true, "group": item})
	})
	g.DELETE("/dns-groups/:id", func(c *gin.Context) { database.DB.Delete(&model.DNSGroup{}, c.Param("id")); c.JSON(http.StatusOK, gin.H{"ok": true}) })

	g.GET("/dns-ddns", func(c *gin.Context) {
		var items []model.DNSDDNS; database.DB.Order("id").Find(&items)
		c.JSON(http.StatusOK, gin.H{"ddns": items})
	})
	g.POST("/dns-ddns", func(c *gin.Context) {
		var item model.DNSDDNS
		if err := c.ShouldBindJSON(&item); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return }
		database.DB.Create(&item); c.JSON(http.StatusCreated, gin.H{"ok": true, "ddns": item})
	})
	g.PUT("/dns-ddns/:id", func(c *gin.Context) {
		var item model.DNSDDNS
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return }
		c.ShouldBindJSON(&item); database.DB.Save(&item); c.JSON(http.StatusOK, gin.H{"ok": true, "ddns": item})
	})
	g.DELETE("/dns-ddns/:id", func(c *gin.Context) { database.DB.Delete(&model.DNSDDNS{}, c.Param("id")); c.JSON(http.StatusOK, gin.H{"ok": true}) })
}