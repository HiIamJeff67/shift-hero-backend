package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	interceptors "github.com/your-org/go-start-monolithic-kit/app/interceptors"
	middlewares "github.com/your-org/go-start-monolithic-kit/app/middlewares"
	modules "github.com/your-org/go-start-monolithic-kit/app/modules"
	metrics "github.com/your-org/go-start-monolithic-kit/app/monitor/metrics"
	constants "github.com/your-org/go-start-monolithic-kit/shared/constants"
)

func configureDevelopmentUserInfoRoutes() {
	userInfoModule := modules.NewUserInfoModule()

	userInfoRoutes := DevelopmentRouterGroup.Group("/userInfo")
	defaultsMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userInfoRoutes.GET(
			"/getMyInfo",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyInfo"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserInfo.GetMyInfo,
					),
				},
				defaultsMiddlewares,
				userInfoModule.Binder.BindGetMyInfo(
					userInfoModule.Controller.GetMyInfo,
				),
			)...,
		)
		userInfoRoutes.PUT(
			"/updateMyInfo",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyInfo"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserInfo.UpdateMyInfo,
					),
				},
				defaultsMiddlewares,
				userInfoModule.Binder.BindUpdateMyInfo(
					userInfoModule.Controller.UpdateMyInfo,
				),
			)...,
		)
	}
}
