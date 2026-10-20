package model

import "time"

type Notification struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Type      string    `gorm:"column:ntype;not null;default:info" json:"type"`
	Title     string    `gorm:"not null" json:"title"`
	Message   string    `gorm:"not null" json:"message"`
	IsRead    bool      `gorm:"column:is_read;default:false" json:"read"`
	CreatedAt time.Time `json:"time"`
}

type Backup struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	BackupType string    `gorm:"column:backup_type;not null;default:manual" json:"type"`
	Size       int64     `gorm:"default:0" json:"size"`
	FileCount  int       `gorm:"column:file_count;default:0" json:"files"`
	Path       *string   `json:"path,omitempty"`
	CreatedAt  time.Time `json:"timestamp"`
}

type AutoBackupSetting struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	Enabled      bool      `gorm:"default:true" json:"enabled"`
	Interval     string    `gorm:"not null;default:daily" json:"interval"`
	Retention    int       `gorm:"default:7" json:"retention"`
	ScheduleTime string    `gorm:"column:schedule_time;not null;default:03:00" json:"time"`
	Path         string    `gorm:"not null;default:/var/backups/openresty" json:"path"`
	UpdatedAt    time.Time
}

type ConfigHistory struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FilePath  string    `gorm:"column:file_path;not null" json:"filepath"`
	Content   string    `gorm:"not null" json:"content"`
	Message   string    `gorm:"default:''" json:"message"`
	SHA256    *string   `json:"sha256,omitempty"`
	CreatedAt time.Time `json:"timestamp"`
}

// ConfigRelease represents one atomic configuration publish operation.
// Status flows through: draft → preparing → writing → validating → reloading → health_checking
// Terminal states: succeeded, validation_failed, reload_failed, rolled_back, recovery_required
type ConfigRelease struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Status      string     `gorm:"not null;default:draft" json:"status"`
	Source      string     `gorm:"not null;default:config" json:"source"`
	Files       string     `gorm:"type:text;default:'[]'" json:"files"`
	Summary     string     `gorm:"default:''" json:"summary"`
	Output      *string    `json:"output,omitempty"`
	Error       *string    `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

type AppSetting struct {
	ID           int64  `gorm:"primaryKey" json:"-"`
	ServerName   string `gorm:"not null;default:OpenResty Admin" json:"serverName"`
	Language     string `gorm:"not null;default:zh_CN" json:"language"`
	Theme        string `gorm:"not null;default:auto" json:"theme"`
	NginxBin     string `gorm:"default:''" json:"nginxBin"`
	NginxConfDir string `gorm:"default:''" json:"nginxConfDir"`
	NginxConfPath string `gorm:"default:''" json:"nginxConfPath"`
	NginxLogDir  string `gorm:"default:''" json:"nginxLogDir"`
	NginxPIDPath string `gorm:"default:''" json:"nginxPIDPath"`
}