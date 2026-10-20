package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/database"
	"openresty-admin/internal/model"
)

var validLuaStages = map[string]bool{
	"init_by_lua_block": true, "init_worker_by_lua_block": true,
	"rewrite_by_lua_block": true, "access_by_lua_block": true,
	"content_by_lua_block": true, "header_filter_by_lua_block": true,
	"body_filter_by_lua_block": true, "log_by_lua_block": true,
	"balancer_by_lua_block": true, "ssl_certificate_by_lua_block": true,
}

var validLuaScopes = map[string]bool{"global": true, "server": true, "location": true}

func RegisterLua(g *gin.RouterGroup) {
	g.GET("/lua", func(c *gin.Context) {
		var items []model.LuaScript
		database.DB.Order("id").Find(&items)
		blocks := make([]gin.H, 0, len(items))
		for _, ls := range items {
			blocks = append(blocks, gin.H{
				"siteName": ls.SiteName, "filePath": ls.FilePath,
				"directive": ls.Directive, "value": ls.Value,
			})
		}
		c.JSON(http.StatusOK, gin.H{"luaBlocks": blocks})
	})
	g.GET("/lua/libraries", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"libraries": []gin.H{
			{"name": "resty.http", "description": "HTTP client for OpenResty", "installed": true},
			{"name": "resty.redis", "description": "Redis client for OpenResty", "installed": true},
			{"name": "resty.mysql", "description": "MySQL client for OpenResty", "installed": true},
			{"name": "resty.websocket", "description": "WebSocket support for OpenResty", "installed": true},
			{"name": "resty.limit.traffic", "description": "Traffic limiting module", "installed": true},
			{"name": "resty.core", "description": "Core OpenResty API", "installed": true},
			{"name": "cjson", "description": "Fast JSON encoding/decoding", "installed": true},
		}})
	})
	g.POST("/lua", func(c *gin.Context) {
		var item model.LuaScript
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
		}
		if !validLuaStages[item.Stage] { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stage: " + item.Stage}); return }
		if !validLuaScopes[item.Scope] { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scope: " + item.Scope}); return }
		database.DB.Create(&item)
		c.JSON(http.StatusCreated, gin.H{"ok": true, "script": item})
	})
	g.PUT("/lua/:id", func(c *gin.Context) {
		var item model.LuaScript
		if err := database.DB.First(&item, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
		}
		c.ShouldBindJSON(&item)
		database.DB.Save(&item)
		c.JSON(http.StatusOK, gin.H{"ok": true, "script": item})
	})
	g.DELETE("/lua/:id", func(c *gin.Context) {
		database.DB.Delete(&model.LuaScript{}, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}