package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"openresty-admin/internal/auth"
	"openresty-admin/internal/model"
)

var DB *gorm.DB

// MustOpen connects to the database using the specified driver ("postgres" or "sqlite").
// dsn is a PostgreSQL connection string or a SQLite file path.
func MustOpen(driver, dsn string) *gorm.DB {
	var dialector gorm.Dialector
	switch driver {
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		dialector = postgres.Open(dsn)
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil { panic("connect database: " + err.Error()) }
	DB = db
	db.AutoMigrate(
		&model.User{}, &model.Session{},
		&model.Site{}, &model.SiteServerName{}, &model.SiteListen{},
		&model.SiteLocation{}, &model.SiteLuaDirective{},
		&model.Stream{},
		&model.Upstream{}, &model.UpstreamServer{},
		&model.Certificate{}, &model.CertificateDomain{}, &model.AcmeUser{},
		&model.AccessList{}, &model.AccessListRule{},
		&model.Snippet{}, &model.Namespace{},
		&model.DNSCredential{}, &model.DNSDomain{}, &model.DNSGroup{}, &model.DNSGroupDomain{}, &model.DNSDDNS{},
		&model.Notification{}, &model.Backup{}, &model.AutoBackupSetting{},
		&model.ConfigHistory{}, &model.ConfigRelease{}, &model.AppSetting{},
		&model.StatsSystem{}, &model.StatsNginx{},
		&model.LuaScript{},
	)
	seed(db)
	return db
}

func seed(db *gorm.DB) {
	adminHash, _ := auth.HashPassword("admin123")
	db.FirstOrCreate(&model.User{}, model.User{Username: "admin", Email: "admin@openresty.local", Password: adminHash, Role: model.RoleAdministrator})
	db.FirstOrCreate(&model.AppSetting{}, model.AppSetting{ID: 1, ServerName: "OpenResty Admin"})
	db.FirstOrCreate(&model.AutoBackupSetting{}, model.AutoBackupSetting{ID: 1})
}