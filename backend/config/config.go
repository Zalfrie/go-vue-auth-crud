package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPass             string
	DBName             string
	JWTSecret          string
	JWTExpirationHours int
	SMTPHost           string
	SMTPPort           int
	SMTPUser           string
	SMTPPass           string
	AppPort            string
	AdminEmail         string
}

// Load reads .env file
func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	return &Config{
		DBHost: viper.GetString("DB_HOST"),
		DBPort: viper.GetString("DB_PORT"),
		DBUser: viper.GetString("DB_USER"),
		DBPass: viper.GetString("DB_PASS"),
		DBName: viper.GetString("DB_NAME"),
		JWTSecret: viper.GetString("JWT_SECRET"),
		JWTExpirationHours: viper.GetInt("JWT_EXPIRATION_HOURS"),
		SMTPHost: viper.GetString("SMTP_HOST"),
		SMTPPort: viper.GetInt("SMTP_PORT"),
		SMTPUser: viper.GetString("SMTP_USER"),
		SMTPPass: viper.GetString("SMTP_PASS"),
		AppPort: viper.GetString("APP_PORT"),
		AdminEmail: viper.GetString("ADMIN_EMAIL"),
	}
}

// ConnectDB opens a DB connection and auto-migrates
func ConnectDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", 
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	return db
}