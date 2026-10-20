package model

import "time"

type StatsSystem struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	CPUCount    int       `gorm:"not null"`
	MemoryUsed  int64     `gorm:"not null"`
	MemoryTotal int64     `gorm:"not null"`
	LoadAvg1    float64   `gorm:"not null"`
	LoadAvg5    float64   `gorm:"not null"`
	LoadAvg15   float64   `gorm:"not null"`
	UptimeSecs  int64     `gorm:"not null"`
	RecordedAt  time.Time `gorm:"not null"`
}

type StatsNginx struct {
	ID                int64     `gorm:"primaryKey;autoIncrement"`
	ActiveConnections int64     `gorm:"not null;default:0"`
	Accepts           int64     `gorm:"not null;default:0"`
	Handled           int64     `gorm:"not null;default:0"`
	Requests          int64     `gorm:"not null;default:0"`
	Reading           int64     `gorm:"not null;default:0"`
	Writing           int64     `gorm:"not null;default:0"`
	Waiting           int64     `gorm:"not null;default:0"`
	RecordedAt        time.Time `gorm:"not null"`
}

type LuaScript struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Version   string    `gorm:"not null;default:1.0.0" json:"version"`
	Stage     string    `gorm:"not null;default:access_by_lua_block" json:"stage"`
	Scope     string    `gorm:"not null;default:server" json:"scope"`
	Directive string    `gorm:"not null" json:"directive"`
	Value     string    `gorm:"not null" json:"value"`
	FilePath  *string   `json:"filePath,omitempty"`
	SiteName  *string   `json:"siteName,omitempty"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time
}