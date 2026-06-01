package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	middlewares "github.com/your-org/go-start-monolithic-kit/app/middlewares"
	constants "github.com/your-org/go-start-monolithic-kit/shared/constants"
)

var (
	DevelopmentRouter      *gin.Engine
	DevelopmentRouterGroup *gin.RouterGroup
)

func ConfigureDevelopmentRoutes() {
	DevelopmentRouterGroup = DevelopmentRouter.Group("/" + constants.DevelopmentBaseURL) // use in development mode
	DevelopmentRouterGroup.Use(
		middlewares.SanitizeXForwardedForMiddleware(),
		middlewares.CORSMiddleware(),
		middlewares.DomainWhiteListMiddleware(),
	)
	DevelopmentRouterGroup.OPTIONS("/*path", func(ctx *gin.Context) { ctx.Status(200) })
	fmt.Println("Router group path:", DevelopmentRouterGroup.BasePath())

	configureDevelopmentAuthRoutes()
	configureDevelopmentUserRoutes()
	configureDevelopmentUserInfoRoutes()
	configureUserSettingRoutes()
	configureDevelopmentUserAccountRoutes()
	configureDevelopmentGraphQLRoutes()

	// test
	configureStaticRoutes()
	configureStorageRoutes()
}
