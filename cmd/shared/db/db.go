// Package db provides a singleton GORM database connection and transaction manager context getter.
// It initializes a PostgreSQL database connection and provides methods to retrieve the database instance
// and transaction manager context getter for use in the application.
package db

import (
	"fmt"
	"os"
	"sync"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"transaction_demo/app/config"
	"transaction_demo/app/constant"
)

var (
	getDBOnce   sync.Once
	dbSingleton *gorm.DB

	getTrmTxGetter sync.Once
	trmTxGetter    *trmgorm.CtxGetter
)

// GetDB returns a singleton instance of GORM database connection.
//
// Returns:
//   - *gorm.DB: Singleton database instance
func GetDB(cf *config.Config) *gorm.DB {
	var err error
	if dbSingleton == nil {
		getDBOnce.Do(func() {
			dbSingleton, err = initDBConnection(cf.Postgres)
			if err != nil {
				os.Exit(constant.ApplicationLoadFailed)
			}
		})
	}
	return dbSingleton
}

// GetTrmGormCtxGetter returns a singleton instance of transaction manager context getter.
// This getter is used to retrieve database transactions from context in transaction-aware operations.
// It uses the default context getter from the go-transaction-manager library.
//
// Returns:
//   - *trmgorm.CtxGetter: Transaction manager context getter instance
func GetTrmGormCtxGetter() *trmgorm.CtxGetter {
	if trmTxGetter == nil {
		getTrmTxGetter.Do(func() {
			trmTxGetter = trmgorm.DefaultCtxGetter
		})
	}
	return trmTxGetter
}

// initDBConnection initializes a new GORM database connection to PostgreSQL.
// Parameters:
//   - cfg: PostgreSQL configuration containing DSN and connection pool settings
//
// Returns:
//   - *gorm.DB: Initialized GORM database instance
//   - error: Error if connection initialization fails
func initDBConnection(cfg config.Postgres) (*gorm.DB, error) {
	var db *gorm.DB
	logMode := logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: cfg.Conn(),
		}),
		&gorm.Config{
			Logger: logMode,
		},
	)
	if err != nil {
		fmt.Println("creating connection to DB failed", "error", err)
		return db, err
	}
	gormer, err := db.DB()
	if err != nil {
		fmt.Println("creating connection to DB failed", "error", err)
		return db, err
	}
	gormer.SetMaxOpenConns(cfg.MaxOpenConns)
	gormer.SetMaxIdleConns(cfg.MaxIdleConns)

	fmt.Println("connected to DB", "dsn", cfg.Conn())
	return db, nil
}
