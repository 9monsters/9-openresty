package model

import "time"

type Upstream struct {
	ID        int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string           `gorm:"not null" json:"name"`
	Strategy  string           `gorm:"not null;default:round_robin" json:"strategy"`
	NodeID    *int64           `json:"-"`
	Servers   []UpstreamServer `gorm:"foreignKey:UpstreamID" json:"servers"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpstreamServer struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UpstreamID  int64  `gorm:"not null;index" json:"-"`
	Address     string `gorm:"not null" json:"address"`
	Port        int    `gorm:"not null;default:80" json:"port"`
	Weight      int    `gorm:"not null;default:1" json:"weight"`
	MaxFails    int    `gorm:"not null;default:3" json:"maxFails"`
	FailTimeout string `gorm:"not null;default:30s" json:"failTimeout"`
	Status      string `gorm:"not null;default:up" json:"status"`
	Seq         int    `gorm:"default:0"`
}

type Certificate struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	Domains   []string   `gorm:"type:text[];default:'{}'" json:"domains"`
	Issuer    string     `gorm:"not null;default:Manual" json:"issuer"`
	Status    string     `gorm:"not null;default:valid" json:"status"`
	AutoRenew bool       `gorm:"default:false" json:"autoRenew"`
	CertPEM   string     `gorm:"default:''" json:"certPEM,omitempty"`
	KeyPEM    string     `gorm:"default:''" json:"keyPEM,omitempty"`
	NotBefore *time.Time `json:"notBefore,omitempty"`
	NotAfter  *time.Time `json:"notAfter,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AcmeUser struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	Email     string `gorm:"not null" json:"email"`
	Provider  string `gorm:"not null;default:Let's Encrypt" json:"provider"`
	ServerURL string `gorm:"column:server_url;not null;default:https://acme-v02.api.letsencrypt.org/directory" json:"server"`
	Status    string `gorm:"not null;default:active" json:"status"`
	CertCount int    `gorm:"default:0" json:"certCount"`
	CreatedAt time.Time
	UpdatedAt time.Time
}