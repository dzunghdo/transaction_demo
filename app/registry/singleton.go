package registry

import (
	"go.uber.org/fx"

	"transaction_demo/app/config"
	"transaction_demo/app/interface/api/route"
	"transaction_demo/cmd/shared/db"
)

// ProvideSingletons provides the singleton instances for DI
var ProvideSingletons = fx.Provide(
	config.InitConfig,
	route.GetEngine,
	db.GetDB,
	db.GetTrmGormCtxGetter,
	db.GetTxManager,
)
