package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"openresty-admin/internal/model"
	"openresty-admin/internal/nginx"
	"openresty-admin/internal/release"
	"openresty-admin/internal/stats"
)

func RegisterDashboard(g *gin.RouterGroup, host *nginx.Host) {
	g.GET("/dashboard", func(c *gin.Context) {
		ver, isOR := host.Version()
		c.JSON(http.StatusOK, gin.H{"status": host.Status(), "version": ver, "isOpenResty": isOR, "stats": host.Stats(), "serverTime": time.Now().UTC().Format(time.RFC3339Nano)})
	})
}

func RegisterSystem(g *gin.RouterGroup, host *nginx.Host) {
	g.GET("/system/info", systemInfo(host))
	g.GET("/system/self-check", selfCheck(host))
	g.GET("/system/upgrade", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"currentVersion": "0.1.0", "latestVersion": "0.1.0", "updateAvailable": false}) })
	g.GET("/system/about", func(c *gin.Context) {
		ver, _ := host.Version()
		c.JSON(http.StatusOK, gin.H{"name": "OpenResty Admin", "version": "0.1.0", "openresty": ver, "goVersion": "go1.25", "dbType": "PostgreSQL", "platform": "darwin", "arch": "amd64"})
	})
}

func RegisterNginx(g *gin.RouterGroup, host *nginx.Host) {
	g.GET("/nginx/status", func(c *gin.Context) {
		ver, isOR := host.Version()
		c.JSON(http.StatusOK, gin.H{"status": host.Status(), "version": ver, "isOpenResty": isOR})
	})
	g.POST("/nginx/test", func(c *gin.Context) { ok, o := host.Test(); c.JSON(http.StatusOK, gin.H{"ok": ok, "output": o}) })
	g.POST("/nginx/reload", func(c *gin.Context) {
		if ok, o := host.Test(); !ok { c.JSON(http.StatusOK, gin.H{"ok": false, "output": o}); return }
		ok, o := host.Reload(); c.JSON(http.StatusOK, gin.H{"ok": ok, "output": o})
	})
	g.POST("/nginx/configure-status", func(c *gin.Context) {
		ch, err := host.StubStatusChanges()
		if err != nil { c.JSON(http.StatusOK, gin.H{"ok": false, "error": err.Error()}); return }
		_ = ch; c.JSON(http.StatusOK, gin.H{"ok": true})
	})

}

func RegisterConfig(g *gin.RouterGroup, host *nginx.Host, db *gorm.DB, coord *release.Coordinator) {
	// Read-only endpoints — host adapter directly.
	g.GET("/config", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"files": host.ListConfigFiles(c.Query("path"))}) })
	g.GET("/config/read", func(c *gin.Context) {
		p := c.Query("path")
		if p == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "path required"}); return }
		content, ok := host.ReadConfig(p)
		if !ok { c.JSON(http.StatusNotFound, gin.H{"error": "File not found"}); return }
		c.JSON(http.StatusOK, gin.H{"path": p, "content": content})
	})

	// Write endpoints — all routed through the release coordinator.
	g.POST("/config", func(c *gin.Context) {
		var body struct{ Path, Content string; Overwrite bool }
		if err := c.ShouldBindJSON(&body); err != nil || body.Path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path and content required"}); return
		}
		if _, ok := host.ReadConfig(body.Path); ok && !body.Overwrite {
			c.JSON(http.StatusConflict, gin.H{"error": "File already exists"}); return
		}
		rel, err := coord.Publish([]nginx.Change{{Path: body.Path, Content: body.Content}}, "config")
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "releaseId": rel.ID, "status": rel.Status})
			return
		}
		db.Create(&model.ConfigHistory{FilePath: body.Path, Content: body.Content, Message: "manual create"})
		c.JSON(http.StatusOK, gin.H{"ok": true, "path": body.Path, "releaseId": rel.ID})
	})
	g.PUT("/config", func(c *gin.Context) {
		var body struct{ Path, Content string }
		if err := c.ShouldBindJSON(&body); err != nil || body.Path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path and content required"}); return
		}
		rel, err := coord.Publish([]nginx.Change{{Path: body.Path, Content: body.Content}}, "config")
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "releaseId": rel.ID, "status": rel.Status})
			return
		}
		db.Create(&model.ConfigHistory{FilePath: body.Path, Content: body.Content, Message: "manual edit"})
		c.JSON(http.StatusOK, gin.H{"ok": true, "path": body.Path, "releaseId": rel.ID})
	})
	g.DELETE("/config", func(c *gin.Context) {
		p := c.Query("path")
		if p == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "path required"}); return }
		rel, err := coord.Publish([]nginx.Change{{Path: p, Delete: true}}, "config")
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "releaseId": rel.ID, "status": rel.Status})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "releaseId": rel.ID})
	})

	// Release history.
	g.GET("/config/releases", func(c *gin.Context) {
		var rels []model.ConfigRelease
		db.Order("created_at DESC").Limit(50).Find(&rels)
		c.JSON(http.StatusOK, gin.H{"releases": rels})
	})
	g.GET("/config/releases/:id", func(c *gin.Context) {
		var rel model.ConfigRelease
		if err := db.First(&rel, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Release not found"}); return
		}
		c.JSON(http.StatusOK, gin.H{"release": rel})
	})
	g.GET("/config/history", func(c *gin.Context) {
		p := c.Query("path")
		if p == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "path required"}); return }
		var entries []model.ConfigHistory
		db.Where("file_path = ?", p).Order("created_at DESC").Find(&entries)
		out := make([]gin.H, 0, len(entries))
		for _, e := range entries { out = append(out, gin.H{"id": e.ID, "filepath": e.FilePath, "timestamp": e.CreatedAt, "message": e.Message}) }
		c.JSON(http.StatusOK, gin.H{"backups": out})
	})
	g.GET("/config/history/:id", func(c *gin.Context) {
		var e model.ConfigHistory
		db.Where("id = ? AND file_path = ?", c.Param("id"), c.Query("path")).First(&e)
		c.JSON(http.StatusOK, gin.H{"path": c.Query("path"), "content": e.Content, "timestamp": e.CreatedAt, "message": e.Message})
	})

	// Release coordinator status — for use by process control (Phase 3).
	g.GET("/config/release-lock", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"locked": coord.IsLocked()})
	})
}

func RegisterLogs(g *gin.RouterGroup, host *nginx.Host) {
	g.GET("/logs/access", func(c *gin.Context) { content := host.ReadLog("access", qint(c, "lines", 200)); c.JSON(http.StatusOK, gin.H{"type": "access", "lines": strings.Count(content, "\n") + 1, "content": content}) })
	g.GET("/logs/error", func(c *gin.Context) { content := host.ReadLog("error", qint(c, "lines", 200)); c.JSON(http.StatusOK, gin.H{"type": "error", "lines": strings.Count(content, "\n") + 1, "content": content}) })
	g.GET("/logs/files", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"files": host.LogFiles()}) })
	g.GET("/logs/file", func(c *gin.Context) {
		p := c.Query("path")
		if p == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "path required"}); return }
		content, err := host.ReadLogByPath(p, qint(c, "lines", 200))
		if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}); return }
		c.JSON(http.StatusOK, gin.H{"content": content})
	})
}

func RegisterStats(g *gin.RouterGroup) {
	g.GET("/stats/system", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"history": stats.SystemHistory.Snapshot()}) })
	g.GET("/stats/nginx", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"history": stats.NginxHistory.Snapshot()}) })
}

func qint(c *gin.Context, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" { return fallback }
	v, _ := strconv.Atoi(raw)
	return v
}