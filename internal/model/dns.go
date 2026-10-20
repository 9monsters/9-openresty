package model

import "time"

type DNSCredential struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	Provider  string `gorm:"not null" json:"provider"`
	APIKey    string `gorm:"column:api_key;not null" json:"apiKey"`
	SecretKey string `gorm:"column:secret_key;default:''" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DNSDomain struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Domain       string `gorm:"not null" json:"domain"`
	CredentialID *int64 `json:"credentialId,omitempty"`
	RecordCount  int    `gorm:"default:0" json:"recordCount"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type DNSGroup struct {
	ID        int64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string   `gorm:"not null" json:"name"`
	Domains   []string `gorm:"type:text[];default:'{}'" json:"domains"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DNSDDNS struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Hostname    string     `gorm:"not null" json:"hostname"`
	Domain      string     `gorm:"not null" json:"domain"`
	IP          *string    `json:"ip,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	CreatedAt   time.Time
}