package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"transaction_demo/app/config"
	"transaction_demo/app/interface/api/handler"
	"transaction_demo/app/interface/api/route"
	"transaction_demo/app/registry"
)

// main initializes the application using Uber Fx framework.
// It sets up the dependency injection, configures the server, and starts listening for requests.
// It provides the necessary components such as repositories, use cases, and handlers.
func main() {
	fx.New(
		registry.ProvideSingletons,
		registry.ProvideRepositories,
		registry.ProvideUsecases,
		fx.Provide(handler.NewAccountHandler),
		fx.Invoke(route.RegisterAccountRoutes),
		fx.Invoke(startServer),
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ConsoleLogger{W: os.Stdout}
		}),
	).Run()
}

// startServer sets up the Gin engine and starts the HTTP server.
func startServer(
	lc fx.Lifecycle,
	engine *gin.Engine,
	cf *config.Config,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := engine.Run(":" + strconv.Itoa(int(cf.Server.Port))); err != nil {
					fmt.Println("start server fail", "error", err)
					panic(err)
				}
			}()
			fmt.Println("start server", "port", cf.Server.Port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("stop server", "port", cf.Server.Port)
			return nil
		},
	})
}
