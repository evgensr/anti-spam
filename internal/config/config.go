package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		GoogleCSV                   string `yaml:"GoogleCSV" envconfig:"APP_GOOGLE_CSV"`
		GoogleDoc                   string `yaml:"GoogleDoc" envconfig:"APP_GOOGLE_DOC"`
		LocalNameCSV                string `yaml:"LocalNameCSV" envconfig:"APP_LOCAL_NAME_CSV"`
		RefreshIntervalSec          int64  `yaml:"RefreshIntervalSec" envconfig:"APP_REFRESH_INTERVAL_SEC"`
		EnableSoundNotification     bool   `yaml:"EnableSoundNotification" envconfig:"APP_ENABLE_SOUND_NOTIFICATION" default:"false"`
		DeleteNotificationNewMember bool   `yaml:"DeleteNotificationNewMember" envconfig:"APP_DELETE_NOTIFICATION_NEW_MEMBER" default:"false"`
	}
	Log struct {
		Path       string `yaml:"Path" envconfig:"LOG_PATH"`
		MaxSize    int    `yaml:"MaxSize" envconfig:"LOG_MAX_SIZE"`
		MaxBackups int    `yaml:"MaxBackups" envconfig:"LOG_MAX_BACKUPS"`
		MaxAge     int    `yaml:"MaxAge" envconfig:"LOG_MAX_AGE"`
		Compress   bool   `yaml:"Compress" envconfig:"LOG_COMPRESS"`
		LocalTime  bool   `yaml:"LocalTime" envconfig:"LOG_LOCAL_TIME"`
	}
	Telegram struct {
		Token   string  `yaml:"Token" envconfig:"TELEGRAM_TOKEN"`
		GroupID int64   `yaml:"GroupID" envconfig:"TELEGRAM_GROUP_ID"`
		AdminID []int64 `yaml:"AdminID" envconfig:"TELEGRAM_ADMIN_ID"`
	}
}

func NewConfig() (*Config, error) {

	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return config, nil
}
