package route

import (
	"transaction_demo/app/interface/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(router *gin.Engine, accountHdl *handler.AccountHandler) {
	apiGroup := router.Group("/api/v1")

	accountGroup := apiGroup.Group("/accounts")
	{
		accountGroup.GET("/:account_id", accountHdl.GetAccountBalance)
		accountGroup.POST("", accountHdl.CreateAccount)
	}

	txGroup := apiGroup.Group("/transactions")
	{
		txGroup.POST("/", accountHdl.MakeTransaction)
	}
}
