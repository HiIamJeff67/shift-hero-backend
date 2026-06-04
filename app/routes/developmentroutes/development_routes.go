package developmentroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	middlewares "github.com/HiIamJeff67/shift-hero-backend/app/middlewares"
	constants "github.com/HiIamJeff67/shift-hero-backend/shared/constants"
)

var (
	DevelopmentRouter      *gin.Engine
	DevelopmentRouterGroup *gin.RouterGroup
)

func ConfigureDevelopmentRoutes() {
	configureHealthRoutes()

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
	configureDevelopmentCompanyRoutes()
	configureDevelopmentSchedulingRoutes()

	// test
	configureStaticRoutes()
	configureStorageRoutes()
}
