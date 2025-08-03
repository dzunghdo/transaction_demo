package registry

import (
	"transaction_demo/app/usecase"

	"go.uber.org/fx"
)

// ProvideUsecases provides the usecase instances for DI
var ProvideUsecases = fx.Provide(usecase.NewAccountUsecase)
