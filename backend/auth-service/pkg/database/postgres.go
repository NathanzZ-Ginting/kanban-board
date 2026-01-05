package database

import (
"fmt"
"log"
"os"
"time"

"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/config"
"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/internal/domain"

"gorm.io/driver/postgres"
"gorm.io/gorm"
"gorm.io/gorm/logger"
)

type PostgresDatabase struct {
	DB *gorm.DB
}

func NewPostgresDatabase(cfg *config.Config) (*PostgresDatabase, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
cfg.DBHost,
cfg.DBUser,
cfg.DBPassword,
cfg.DBName,
cfg.DBPort,
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Supabase PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PostgresDatabase{DB: db}, nil
}

func (p *PostgresDatabase) AutoMigrate() error {
	return p.DB.AutoMigrate(
&domain.User{},
		&domain.RefreshToken{},
	)
}

func (p *PostgresDatabase) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgresDatabase) GetDB() *gorm.DB {
	return p.DB
}
