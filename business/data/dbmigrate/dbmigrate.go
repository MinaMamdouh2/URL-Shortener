// Package dbmigrate contains the database schema, migrations and seeding data.Add commentMore actions
// We will use darwin package, it is simple and what it does it lets the project maintain the schema, update the schema
// in a programmable way, so no extra cli tooling.
// What darwin is going to do is let us define a file with a certain structure that maintains the changes we wanna migrate
// to the DB, it adds a table to a DB that records what has been done already so we can run it over and over and not
// having any problems. It does create a hash of the code that is applying, so if somebody changes something things will
// fail because the idea is you run it once and you don't change it.
// There are some problems here, once we update or migrate the DB if we still have services that is running against old
// schema those calls are going to fail "chicken and egg problem between the time you update the DB and from the time you
// update the services" if you are in a simple environment "local environment", ideally as we update our service we will
// do the migrations on the DB first and the system is kinda down during that transition but most people don't want to be
// down during a maintenance or an upgrade, you gotta find a way.
package dbmigrate

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/ardanlabs/darwin/v3"
	"github.com/ardanlabs/darwin/v3/dialects/postgres"
	"github.com/ardanlabs/darwin/v3/drivers/generic"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

// We embed the sql files and this sql will be part of the admin binary, we are going to use to apply this migration
// and it will become part of the project.
var (
	//go:embed sql/migrate.sql
	migrateDoc string

	//go:embed sql/seed.sql
	seedDoc string
)

// Migrate attempts to bring the database up to date with the migrations
// defined in this package.
// The migrate call is assuming that we got a DB connection using sqlx and that our DB package has a sort of function
// that the DB is available
func Migrate(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting raw database handle: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	// here we are using darwin level APIs, we create a new postgres dialect which means we are using postgres syntax to do
	// the things we want to do.
	driver, err := generic.New(sqlDB, postgres.Dialect{})
	if err != nil {
		return fmt.Errorf("construct darwin driver: %w", err)
	}
	migrations := darwin.ParseMigrations(migrateDoc)
	fmt.Printf("Parsed %d migrations\n", len(migrations))
	for _, m := range migrations {
		fmt.Printf("Version: %f, Description: %s\n", m.Version, m.Description)
	}
	// Here ask darwin to read the migration file, parse it into a set of objects that "New" knows how to read against
	// the driver we are going to be using for the DB
	d := darwin.New(driver, darwin.ParseMigrations(migrateDoc))
	// Apply the migrations to our DB.
	return d.Migrate()
}

// Seed runs the seed document defined in this package against db. The queries
// are run in a transaction and rolled back if any fail.
func Seed(ctx context.Context, db *gorm.DB) (err error) {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting raw database handle: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	// Start a transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// We do defer rollback, if a commit is already happened then the rollback is ignored, but if the commit hasn't happened
	// then we will end up being able to rollback
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback().Error; rbErr != nil {
				if errors.Is(rbErr, sql.ErrTxDone) {
					return
				}
				err = fmt.Errorf("rollback error: %w", rbErr)
			}
		}

	}()
	// Execute everything in the seed doc
	if err = tx.Exec(seedDoc).Error; err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	// Commit it
	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// SeedCustom runs the specified seed document against db. The queries are run
// in a transaction and rolled back if any fail.
func SeedCustom(ctx context.Context, db *gorm.DB, seedDoc string) (err error) {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting raw database handle: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	// Start a transaction
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if err != nil {
			if errTx := tx.Rollback().Error; errTx != nil {
				if errors.Is(errTx, sql.ErrTxDone) {
					return
				}
				err = fmt.Errorf("rollback: %w", errTx)
				return
			}
		}
	}()

	if err = tx.Exec(seedDoc).Error; err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
