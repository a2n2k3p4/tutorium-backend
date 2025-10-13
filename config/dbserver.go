package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func (c *Config) DBUrl() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
}

func ConnectDB(cfg *Config) (*gorm.DB, error) {
	dbUrl := cfg.DBUrl()
	dbConfig := &gorm.Config{}

	switch GORMLog() {
	case "silent":
		dbConfig.Logger = logger.Default.LogMode(logger.Silent)
	case "error":
		dbConfig.Logger = logger.Default.LogMode(logger.Error)
	case "warn":
		dbConfig.Logger = logger.Default.LogMode(logger.Warn)
	case "info":
		dbConfig.Logger = logger.Default.LogMode(logger.Info)
	default:
		dbConfig.Logger = logger.Default.LogMode(logger.Warn)
	}

	db, err := gorm.Open(postgres.Open(dbUrl), dbConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

func NewConfig() *Config {
	return &Config{
		DBUser:     DBUser(),
		DBPassword: DBPassword(),
		DBHost:     DBHost(),
		DBPort:     DBPort(),
		DBName:     DBName(),
	}
}
