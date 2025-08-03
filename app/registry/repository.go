package registry

import (
	"transaction_demo/app/external/persist/postgres"

	"go.uber.org/fx"
)

// ProvideRepositories provides the repository instances for DI
var ProvideRepositories = fx.Provide(postgres.NewAccountRepository, postgres.NewTransactionRepository)
