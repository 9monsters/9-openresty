package model

import "time"

type AccessList struct {
	ID          int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string           `gorm:"not null" json:"name"`
	Description string           `gorm:"default:''" json:"description"`
	Rules       []AccessListRule `gorm:"foreignKey:ListID" json:"rules"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AccessListRule struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ListID   int64  `gorm:"not null;index" json:"-"`
	RuleType string `gorm:"not null;default:allow" json:"type"`
	IP       string `gorm:"not null" json:"ip"`
	CIDR     *int   `json:"cidr,omitempty"`
	Seq      int    `gorm:"default:0"`
}

type Snippet struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Content     string    `gorm:"default:''" json:"content"`
	Description string    `gorm:"default:''" json:"description"`
	IsTemplate  bool      `gorm:"default:false" json:"isTemplate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time
}

type Namespace struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"not null;uniqueIndex" json:"name"`
	Description string    `gorm:"default:''" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}