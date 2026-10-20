package model

import "time"

// Role constants define the three permission levels.
const (
	RoleViewer        = "Viewer"
	RoleOperator      = "Operator"
	RoleAdministrator = "Administrator"
)

// RoleLevel maps a role name to its numeric privilege level.
var RoleLevel = map[string]int{
	RoleViewer:        1,
	RoleOperator:      2,
	RoleAdministrator: 3,
}

// RequireRole returns true when the user's role meets or exceeds the required level.
func RequireRole(userRole, required string) bool {
	return RoleLevel[userRole] >= RoleLevel[required]
}

type User struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string     `gorm:"not null;unique" json:"username"`
	Email     string     `gorm:"not null"                json:"email"`
	Password  string     `gorm:"not null"                json:"-"`
	Role      string     `gorm:"not null;default:Operator" json:"role"`
	Status    string     `gorm:"not null;default:active" json:"status"`
	TwoFactor bool       `gorm:"not null;default:false"  json:"twoFactorEnabled"`
	Theme     string     `gorm:"not null;default:auto"   json:"theme"`
	Language  string     `gorm:"not null;default:zh_CN"  json:"language"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type Session struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index"`
	TokenJTI  string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}