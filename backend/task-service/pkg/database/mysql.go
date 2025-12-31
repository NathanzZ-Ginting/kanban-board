package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/config"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDatabase struct {
	DB *gorm.DB
}

func NewMySQLDatabase(cfg *config.Config) (*MySQLDatabase, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDatabase,
	)

	var gormLogger logger.Interface
	if cfg.Environment == "development" {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &MySQLDatabase{DB: db}, nil
}

func (m *MySQLDatabase) AutoMigrate() error {
	return m.DB.AutoMigrate(
		&domain.Task{},
		&domain.Comment{},
		&domain.Label{},
	)
}

func (m *MySQLDatabase) Close() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (m *MySQLDatabase) GetDB() *gorm.DB {
	return m.DB
}
