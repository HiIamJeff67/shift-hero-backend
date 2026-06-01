package graphql

import (
	"github.com/gin-gonic/gin"
)

func GraphQLHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{
			"message": "GraphQL is disabled in this template baseline.",
		})
	}
}

func PlaygroundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{
			"message": "GraphQL playground is disabled in this template baseline.",
		})
	}
}
