package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/nginx"
	"openresty-admin/internal/release"
)

func RegisterProcess(g *gin.RouterGroup, host *nginx.Host, coord *release.Coordinator) {
	g.POST("/nginx/stop", func(c *gin.Context) {
		if coord.IsLocked() {
			c.JSON(http.StatusConflict, gin.H{"ok": false, "error": "A configuration release is in progress. Wait for it to complete before stopping nginx."})
			return
		}
		if ok, out := host.Test(); !ok {
			c.JSON(http.StatusPreconditionFailed, gin.H{"ok": false, "error": "nginx -t failed", "output": out})
			return
		}
		ok, out := host.Stop()
		c.JSON(http.StatusOK, gin.H{"ok": ok, "output": out, "status": host.Status()})
	})

	g.POST("/nginx/start", func(c *gin.Context) {
		if coord.IsLocked() {
			c.JSON(http.StatusConflict, gin.H{"ok": false, "error": "A configuration release is in progress."})
			return
		}
		ok, out := host.Start()
		c.JSON(http.StatusOK, gin.H{"ok": ok, "output": out, "status": host.Status()})
	})

	g.POST("/nginx/restart", func(c *gin.Context) {
		if coord.IsLocked() {
			c.JSON(http.StatusConflict, gin.H{"ok": false, "error": "A configuration release is in progress."})
			return
		}
		if ok, out := host.Test(); !ok {
			c.JSON(http.StatusPreconditionFailed, gin.H{"ok": false, "error": "nginx -t failed", "output": out})
			return
		}
		stopOK, stopOut := host.Stop()
		if !stopOK {
			c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "stop failed", "output": stopOut})
			return
		}
		startOK, startOut := host.Start()
		c.JSON(http.StatusOK, gin.H{"ok": startOK, "output": stopOut + "\n" + startOut, "status": host.Status()})
	})
}
