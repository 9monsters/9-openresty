package model

import "time"


type Site struct {
	ID            int64              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string             `gorm:"not null" json:"name"`
	FilePath      string             `gorm:"not null" json:"filePath"`
	ServerNames   []string           `gorm:"type:text[];default:'{}'" json:"serverName"`
	Listens       []string           `gorm:"type:text[];default:'{80}'" json:"listen"`
	SSL           bool               `gorm:"default:false" json:"ssl"`
	HTTP2         bool               `gorm:"default:false" json:"http2"`
	RootDir       *string            `json:"root,omitempty"`
	AccessLog     *string            `json:"accessLog,omitempty"`
	ErrorLog      *string            `json:"errorLog,omitempty"`
	ProxyPass     *string            `json:"proxyPass,omitempty"`
	Content       string             `gorm:"not null;default:''" json:"content"`
	Enabled       bool               `gorm:"default:true" json:"enabled"`
	Locations     []SiteLocation     `gorm:"foreignKey:SiteID" json:"locations"`
	LuaDirectives []SiteLuaDirective `gorm:"foreignKey:SiteID" json:"luaDirectives"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time
}

type SiteLocation struct {
	ID              int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SiteID          int64  `gorm:"not null;index" json:"-"`
	Path            string `gorm:"not null;default:/" json:"path"`
	ProxyPass       *string `json:"proxyPass,omitempty"`
	RootDir         *string `json:"root,omitempty"`
	ReturnDirective *string `json:"return,omitempty"`
	Seq             int    `gorm:"default:0"`
}

type SiteLuaDirective struct {
	ID         int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SiteID     *int64 `json:"-"`
	LocationID *int64 `json:"-"`
	Directive  string `gorm:"not null" json:"type"`
	Value      string `gorm:"not null" json:"value"`
}

type Stream struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Listen    string    `gorm:"column:listen_addr;not null" json:"listen"`
	ProxyPass *string    `json:"proxyPass,omitempty"`
	Protocol  string    `gorm:"not null;default:tcp" json:"protocol"`
	SSL       bool      `gorm:"default:false" json:"ssl"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time
	UpdatedAt time.Time
}