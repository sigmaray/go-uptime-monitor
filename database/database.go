package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"time"

	"go-uptime-monitor/config"

	"github.com/pressly/goose/v3"
	gooseLock "github.com/pressly/goose/v3/lock"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DSN(c config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}

func Connect(cfg config.DatabaseConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(DSN(cfg)), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db
}

func runMigrations(migrations embed.FS, sqlDB *sql.DB) {
	migrationFS, err := fs.Sub(migrations, "migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open migrations directory")
	}

	sessionLocker, err := gooseLock.NewPostgresSessionLocker()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migration session locker")
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		sqlDB,
		migrationFS,
		goose.WithSessionLocker(sessionLocker),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migration provider")
	}

	if _, err := provider.Up(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
}

func RunMigrations(migrations embed.FS, cfg config.DatabaseConfig) {
	db := Connect(cfg)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database handle")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database")
		}
	}()

	runMigrations(migrations, sqlDB)
}
