package developmentroutes

import "github.com/gin-gonic/gin"

func configureHealthRoutes() {
	healthHandler := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	}
	headHealthHandler := func(ctx *gin.Context) {
		ctx.Status(200)
	}

	DevelopmentRouter.GET("/", healthHandler)
	DevelopmentRouter.HEAD("/", headHealthHandler)
	DevelopmentRouter.GET("/healthz", healthHandler)
	DevelopmentRouter.HEAD("/healthz", headHealthHandler)
}
