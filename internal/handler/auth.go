package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"openresty-admin/internal/auth"
	"openresty-admin/internal/middleware"
	"openresty-admin/internal/model"
)

func RegisterAuth(g *gin.RouterGroup, db *gorm.DB) {
	// Login/Refresh are registered separately in router.go at root level
}

func Login(c *gin.Context, db *gorm.DB) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"}); return
	}
	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"}); return
	}
	var user model.User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}); return
	}
	if !auth.CheckPassword(user.Password, body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}); return
	}
	claims := auth.NewClaims(user.ID, user.Username, user.Role)
	token, err := auth.Sign(claims)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "sign token: " + err.Error()}); return }
	db.Create(&model.Session{UserID: user.ID, TokenJTI: claims.Jti, ExpiresAt: time.Unix(claims.Exp, 0)})
	db.Model(&user).Update("last_login", time.Now())
	c.JSON(http.StatusOK, gin.H{"token": token, "expiresAt": time.Unix(claims.Exp, 0).UTC().Format(time.RFC3339),
		"user": gin.H{"id": user.ID, "username": user.Username, "role": user.Role}})
}

func Refresh(c *gin.Context) {
	claims := middleware.Claims(c)
	if claims == nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "no session"}); return }
	nc := auth.NewClaims(claims.UserID, claims.Username, claims.Role)
	token, _ := auth.Sign(nc)
	c.JSON(http.StatusOK, gin.H{"token": token, "expiresAt": time.Unix(nc.Exp, 0).UTC().Format(time.RFC3339)})
}