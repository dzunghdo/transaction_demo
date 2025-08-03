package db

import (
	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"gorm.io/gorm"
)

// GetTxManager returns a transaction manager for the given GORM database instance.
// It uses the default transaction manager factory for GORM and sets the propagation to Nested.
func GetTxManager(db *gorm.DB) trm.Manager {
	return manager.Must(
		trmgorm.NewDefaultFactory(db),
		manager.WithSettings(trmgorm.MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationNested))),
		),
	)
}
