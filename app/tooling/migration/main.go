package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MinaMamdouh2/URL-Shortener/business/data/dbmigrate"
	"github.com/ardanlabs/conf/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if err := migrateSeed(); err != nil {
		log.Fatalln(err)
	}
}

func migrateSeed() error {
	var cfg struct {
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:admin,mask"`
			Host         string `conf:"default:localhost"`
			Port         int    `conf:"default:5432"`
			Name         string `conf:"default:url-shortener"`
			MaxIdleConns int    `conf:"default:2"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
	}

	const prefix = "URL_SHORTENER"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		return fmt.Errorf("parsing config: %w", err)
	}

	// 2) Build GORM DSN (Postgres)
	sslMode := "require"
	if cfg.DB.DisableTLS {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		sslMode,
	)

	sqlConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open postgres connection: %w", err)
	}
	// 3) Open GORM with a minimal logger
	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn: sqlConn,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		return fmt.Errorf("connect database (GORM): %w", err)
	}

	// 4) Configure the underlying *sql.DB for pooling & close
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("getting raw DB from GORM: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	// Zero means unlimited – matches sql.DB behavior
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	defer sqlDB.Close()
	// We setout a 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We call migrate
	if err := dbmigrate.Migrate(ctx, gormDB); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	// ========================================================================
	// We call seed with the same connection
	if err := dbmigrate.Seed(ctx, gormDB); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")

	return nil
}
