package route

import (
	"sync"

	"transaction_demo/app/interface/api/middleware"

	"github.com/gin-gonic/gin"
)

var (
	routeOnce sync.Once
	router    *gin.Engine
)

// GetEngine initializes and returns the Gin engine for the application.
// It ensures that the engine is initialized only once and sets up common middleware
// for logging and error recovery.
//
// Returns:
//   - *gin.Engine: The initialized Gin engine
func GetEngine() *gin.Engine {
	if router == nil {
		routeOnce.Do(func() {
			router = gin.New()
			router.Use(
				gin.LoggerWithConfig(gin.LoggerConfig{
					Output: gin.DefaultWriter,
				}),
				middleware.Recover(),
			)
		})
	}

	return router
}
