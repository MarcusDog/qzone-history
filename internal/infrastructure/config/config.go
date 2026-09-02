package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Type     string // 数据库类型，例如 "sqlite", "postgres", "mysql"
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	QzoneAPI struct {
		QRCodeURL  string
		LoginURL   string
		MessageURL string
	}
}

func LoadConfig() (*Config, error) {
	// 设置默认配置
	defaultConfig := &Config{
		Database: struct {
			Type     string
			Host     string
			Port     int
			User     string
			Password string
			DBName   string
			SSLMode  string
		}{
			Type:     "sqlite",
			Host:     "localhost",
			Port:     5432,
			User:     "your_db_user",
			Password: "your_db_password",
			DBName:   "./qzone-history.db",
			SSLMode:  "disable",
		},
		QzoneAPI: struct {
			QRCodeURL  string
			LoginURL   string
			MessageURL string
		}{
			QRCodeURL:  "https://ssl.ptlogin2.qq.com/ptqrshow?appid=549000912&e=2&l=M&s=3&d=72&v=4&t=0.8692955245720428&daid=5",
			LoginURL:   "https://ssl.ptlogin2.qq.com/ptqrlogin",
			MessageURL: "https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all",
		},
	}

	// 使用独立实例，避免全局 Viper 状态在测试或嵌入使用时相互污染。
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.SetEnvPrefix("QZONE_HISTORY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetDefault("database.type", defaultConfig.Database.Type)
	v.SetDefault("database.host", defaultConfig.Database.Host)
	v.SetDefault("database.port", defaultConfig.Database.Port)
	v.SetDefault("database.user", defaultConfig.Database.User)
	v.SetDefault("database.password", defaultConfig.Database.Password)
	v.SetDefault("database.dbName", defaultConfig.Database.DBName)
	v.SetDefault("database.sslMode", defaultConfig.Database.SSLMode)
	v.SetDefault("qzoneAPI.qrCodeURL", defaultConfig.QzoneAPI.QRCodeURL)
	v.SetDefault("qzoneAPI.loginURL", defaultConfig.QzoneAPI.LoginURL)
	v.SetDefault("qzoneAPI.messageURL", defaultConfig.QzoneAPI.MessageURL)

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	// 将读取的配置与默认配置合并
	if err := v.Unmarshal(defaultConfig); err != nil {
		return nil, err
	}

	return defaultConfig, nil
}
