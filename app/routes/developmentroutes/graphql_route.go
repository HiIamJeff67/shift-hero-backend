package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	graphql "github.com/your-org/go-start-monolithic-kit/app/graphql"
	interceptors "github.com/your-org/go-start-monolithic-kit/app/interceptors"
	middlewares "github.com/your-org/go-start-monolithic-kit/app/middlewares"
)

func configureDevelopmentGraphQLRoutes() {
	graphqlRoutes := DevelopmentRouterGroup.Group("/graphql")

	graphqlRoutes.Use(
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	)
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler())
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
