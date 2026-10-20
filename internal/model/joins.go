package model

// SiteServerName stores individual server_name entries for a Site.
// Replaces the PostgreSQL text[] column.
type SiteServerName struct {
	ID     int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SiteID int64  `gorm:"not null;index" json:"-"`
	Name   string `gorm:"not null" json:"name"`
	Seq    int    `gorm:"default:0"`
}

// SiteListen stores individual listen directives for a Site.
type SiteListen struct {
	ID     int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SiteID int64  `gorm:"not null;index" json:"-"`
	Listen string `gorm:"not null" json:"listen"`
	Seq    int    `gorm:"default:0"`
}

// CertificateDomain stores individual domain names for a Certificate.
type CertificateDomain struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CertificateID int64  `gorm:"not null;index" json:"-"`
	Domain        string `gorm:"not null" json:"domain"`
}

// DNSGroupDomain stores individual domain names for a DNSGroup.
type DNSGroupDomain struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID int64  `gorm:"not null;index" json:"-"`
	Domain  string `gorm:"not null" json:"domain"`
}
