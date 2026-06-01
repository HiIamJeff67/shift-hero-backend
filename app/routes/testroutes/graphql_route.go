package testroutes

import (
	"github.com/your-org/go-start-monolithic-kit/app/graphql"
	middlewares "github.com/your-org/go-start-monolithic-kit/app/middlewares"

	"github.com/gin-gonic/gin"
)

func ConfigureTestGraphQLRoutes() {
	graphqlRoutes := TestRouterGroup.Group("/graphql")

	graphqlRoutes.Use(middlewares.AuthMiddleware())
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler())
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
