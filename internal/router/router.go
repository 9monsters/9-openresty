package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"openresty-admin/internal/embedweb"
	"openresty-admin/internal/handler"
	"openresty-admin/internal/metrics"
	"openresty-admin/internal/middleware"
	"openresty-admin/internal/nginx"
	"openresty-admin/internal/release"
)

func Setup(r *gin.Engine, db *gorm.DB, host *nginx.Host, coord *release.Coordinator) {
	r.Use(middleware.Logger(), gin.Recovery())

	api := r.Group("/api")
	api.Use(middleware.Auth())

	// --- Viewer: read-only access (all authenticated users) ---
	handler.RegisterDashboard(api, host)
	handler.RegisterSystem(api, host)
	handler.RegisterLogs(api, host)
	handler.RegisterStats(api)
	handler.RegisterAuth(api, db)

	// --- Operator: configuration editing ---
	oper := api.Group("")
	oper.Use(middleware.RequireOperator())
	handler.RegisterConfig(oper, host, db, coord)
	handler.RegisterUpstreams(oper)
	handler.RegisterStreams(oper)
	handler.RegisterCertificates(oper)
	handler.RegisterAcme(oper)
	handler.RegisterSnippets(oper)
	handler.RegisterAccessLists(oper)
	handler.RegisterNamespaces(oper)
	handler.RegisterDNS(oper)
	handler.RegisterLua(oper)
	handler.RegisterNotifications(oper, db)
	handler.RegisterBackups(oper, db)

	// --- Administrator: users, nginx control, settings, process ---
	admin := api.Group("")
	admin.Use(middleware.RequireAdmin())
	handler.RegisterNginx(admin, host)
	handler.RegisterProcess(admin, host, coord)
	handler.RegisterUsers(admin, db)
	handler.RegisterSettings(admin, host, db)

	r.POST("/api/auth/login", func(c *gin.Context) { handler.Login(c, db) })
	r.POST("/api/auth/refresh", func(c *gin.Context) { handler.Refresh(c) })
	r.GET("/metrics", gin.WrapH(metrics.Handler))

	// Frontend — embedded or external proxy.
	indexHTML, indexErr := embedweb.ReadFile("index.html")
	if indexErr == nil {
		r.StaticFS("/assets", embedweb.AssetsFS())
		serveIndex := func(c *gin.Context) { c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML) }
		r.GET("/", serveIndex)
		r.GET("/index.html", serveIndex)
		r.GET("/favicon.ico", func(c *gin.Context) {
			data, err := embedweb.ReadFile("favicon.ico")
			if err != nil { c.Status(http.StatusNotFound); return }
			c.Data(http.StatusOK, "image/x-icon", data)
		})
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/metrics") {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
		})
	}
}