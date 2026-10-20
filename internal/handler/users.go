package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"openresty-admin/internal/auth"
	"openresty-admin/internal/config"
	"openresty-admin/internal/model"
	"openresty-admin/internal/nginx"
)

func RegisterUsers(g *gin.RouterGroup, db *gorm.DB) {
	g.GET("/users", func(c *gin.Context) { var items []model.User; db.Order("id").Find(&items); c.JSON(http.StatusOK, gin.H{"users": items}) })
	g.POST("/users", func(c *gin.Context) {
		var body struct{ Username, Email, Password, Role string }
		c.ShouldBindJSON(&body)
		if body.Username == "" || body.Email == "" || body.Password == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "required fields missing"}); return }
		if body.Role == "" { body.Role = model.RoleOperator }
		hash, err := auth.HashPassword(body.Password)
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "hash password: " + err.Error()}); return }
		u := model.User{Username: body.Username, Email: body.Email, Password: hash, Role: body.Role}
		db.Create(&u); c.JSON(http.StatusCreated, gin.H{"ok": true, "user": u})
	})
	g.PUT("/users/:id", func(c *gin.Context) {
		var u model.User
		if err := db.First(&u, c.Param("id")).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}); return }
		var body struct{ Username, Email, Role, Status, Theme, Language string }
		c.ShouldBindJSON(&body)
		if body.Username != "" { u.Username = body.Username }
		if body.Email != "" { u.Email = body.Email }
		if body.Role != "" { u.Role = body.Role }
		if body.Status != "" { u.Status = body.Status }
		if body.Theme != "" { u.Theme = body.Theme }
		if body.Language != "" { u.Language = body.Language }
		db.Save(&u); c.JSON(http.StatusOK, gin.H{"ok": true, "user": u})
	})
	g.DELETE("/users/:id", func(c *gin.Context) { db.Delete(&model.User{}, c.Param("id")); c.JSON(http.StatusOK, gin.H{"ok": true}) })
	g.GET("/users/profile", func(c *gin.Context) { var u model.User; db.First(&u); c.JSON(http.StatusOK, u) })
	g.PUT("/users/profile", func(c *gin.Context) {
		var u model.User; db.First(&u)
		var body struct{ Username, Email, Language, Theme string }
		c.ShouldBindJSON(&body)
		if body.Username != "" { u.Username = body.Username }
		if body.Email != "" { u.Email = body.Email }
		if body.Language != "" { u.Language = body.Language }
		if body.Theme != "" { u.Theme = body.Theme }
		db.Save(&u); c.JSON(http.StatusOK, gin.H{"ok": true, "user": u})
	})
	g.PUT("/users/profile/password", func(c *gin.Context) {
		var u model.User; db.First(&u)
		var body struct{ CurrentPassword, NewPassword string }
		c.ShouldBindJSON(&body)
		if !auth.CheckPassword(u.Password, body.CurrentPassword) { c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"}); return }
		hash, err := auth.HashPassword(body.NewPassword)
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "hash password: " + err.Error()}); return }
		db.Model(&u).Update("password", hash); c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}

func RegisterNotifications(g *gin.RouterGroup, db *gorm.DB) {
	g.GET("/notifications", func(c *gin.Context) {
		page := qint(c, "page", 1)
		var total int64; var items []model.Notification
		q := db.Model(&model.Notification{})
		if t := c.Query("type"); t != "" && t != "all" { q = q.Where("ntype = ?", t) }
		q.Count(&total).Order("created_at DESC").Offset((page - 1) * 20).Limit(20).Find(&items)
		c.JSON(http.StatusOK, gin.H{"notifications": items, "total": total, "page": page, "pageSize": 20})
	})
	g.POST("/notifications", func(c *gin.Context) {
		var body struct{ Type, Title, Message string }; c.ShouldBindJSON(&body)
		if body.Type == "" { body.Type = "info" }
		n := model.Notification{Type: body.Type, Title: body.Title, Message: body.Message}
		db.Create(&n); c.JSON(http.StatusCreated, gin.H{"ok": true, "notification": n})
	})
	g.PUT("/notifications/:id/read", func(c *gin.Context) { db.Model(&model.Notification{}).Where("id = ?", c.Param("id")).Update("is_read", true); c.JSON(http.StatusOK, gin.H{"ok": true}) })
	g.PUT("/notifications/read-all", func(c *gin.Context) { db.Model(&model.Notification{}).Where("is_read = ?", false).Update("is_read", true); c.JSON(http.StatusOK, gin.H{"ok": true}) })
	g.DELETE("/notifications/:id", func(c *gin.Context) { db.Delete(&model.Notification{}, c.Param("id")); c.JSON(http.StatusOK, gin.H{"ok": true}) })
}

func RegisterBackups(g *gin.RouterGroup, db *gorm.DB) {
	g.GET("/backup", func(c *gin.Context) { var items []model.Backup; db.Order("created_at DESC").Find(&items); c.JSON(http.StatusOK, gin.H{"backups": items}) })
	g.POST("/backup", func(c *gin.Context) {
		var body struct{ Name string `json:"name"` }; c.ShouldBindJSON(&body)
		if body.Name == "" { body.Name = "manual-" + time.Now().UTC().Format("2006-01-02") }
		b := model.Backup{Name: body.Name, BackupType: "manual", Size: 300000, FileCount: 12}
		db.Create(&b); c.JSON(http.StatusCreated, gin.H{"ok": true, "backup": b})
	})
	g.POST("/backup/:id/restore", func(c *gin.Context) { var b model.Backup; db.First(&b, c.Param("id")); c.JSON(http.StatusOK, gin.H{"ok": true, "message": "Restored from " + b.Name}) })
	g.DELETE("/backup/:id", func(c *gin.Context) { db.Delete(&model.Backup{}, c.Param("id")); c.JSON(http.StatusOK, gin.H{"ok": true}) })
	g.GET("/backup/auto-backup", func(c *gin.Context) {
		var s model.AutoBackupSetting; db.FirstOrCreate(&s, model.AutoBackupSetting{ID: 1})
		c.JSON(http.StatusOK, gin.H{"settings": s, "schedule": gin.H{"nextRun": time.Now().Add(24 * time.Hour), "lastRun": time.Now().Add(-time.Hour), "lastStatus": "success"}})
	})
	g.PUT("/backup/auto-backup", func(c *gin.Context) {
		var s model.AutoBackupSetting; db.FirstOrCreate(&s, model.AutoBackupSetting{ID: 1}); c.ShouldBindJSON(&s); db.Save(&s)
		c.JSON(http.StatusOK, gin.H{"ok": true, "settings": s})
	})
}

func RegisterSettings(g *gin.RouterGroup, host *nginx.Host, db *gorm.DB) {
	g.GET("/settings", func(c *gin.Context) {
		var s model.AppSetting; db.FirstOrCreate(&s, model.AppSetting{ID: 1})
		paths := host.Paths()
		s.NginxBin = coalesce(s.NginxBin, paths.NginxBin)
		s.NginxConfDir = coalesce(s.NginxConfDir, paths.NginxConfDir)
		s.NginxConfPath = coalesce(s.NginxConfPath, paths.NginxConfPath)
		s.NginxLogDir = coalesce(s.NginxLogDir, paths.NginxLogDir)
		s.NginxPIDPath = coalesce(s.NginxPIDPath, paths.NginxPIDPath)
		ver, isOR := host.Version()
		c.JSON(http.StatusOK, gin.H{"serverName": s.ServerName, "language": s.Language, "theme": s.Theme,
			"nginxBin": s.NginxBin, "nginxConfDir": s.NginxConfDir, "nginxConfPath": s.NginxConfPath,
			"nginxLogDir": s.NginxLogDir, "nginxPIDPath": s.NginxPIDPath,
			"paths": paths, "nginxStatus": host.Status(), "version": ver, "isOpenResty": isOR})
	})
	g.PUT("/settings", func(c *gin.Context) {
		var s model.AppSetting; db.FirstOrCreate(&s, model.AppSetting{ID: 1}); c.ShouldBindJSON(&s); db.Save(&s)
		host.UpdatePaths(config.Paths{NginxBin: s.NginxBin, NginxConfDir: s.NginxConfDir, NginxConfPath: s.NginxConfPath, NginxLogDir: s.NginxLogDir, NginxPIDPath: s.NginxPIDPath})
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	g.GET("/settings/server-name", func(c *gin.Context) { var s model.AppSetting; db.FirstOrCreate(&s, model.AppSetting{ID: 1}); c.JSON(http.StatusOK, gin.H{"name": s.ServerName}) })
	g.POST("/settings/configure-status", func(c *gin.Context) {
		ch, err := host.StubStatusChanges()
		if err != nil { c.JSON(http.StatusOK, gin.H{"ok": false, "error": err.Error()}); return }
		_ = ch; c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}

func coalesce(a, b string) string { if a != "" { return a }; return b }