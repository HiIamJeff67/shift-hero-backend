package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	interceptors "github.com/HiIamJeff67/shift-hero-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/shift-hero-backend/app/middlewares"
	modules "github.com/HiIamJeff67/shift-hero-backend/app/modules"
	constants "github.com/HiIamJeff67/shift-hero-backend/shared/constants"
)

func configureDevelopmentCompanyRoutes() {
	module := modules.NewCompanyModule()

	companyRoutes := DevelopmentRouterGroup.Group("/companies")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(2 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}

	companyRoutes.POST(
		"",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createCompany"),
				middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.company.createCompany"),
			},
			defaultMiddlewares,
			module.Binder.BindCreateCompany(module.Controller.CreateCompany),
		)...,
	)
	companyRoutes.GET(
		"/me",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyCompanies"),
				middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.company.getMyCompanies"),
			},
			defaultMiddlewares,
			module.Binder.BindGetMyCompanies(module.Controller.GetMyCompanies),
		)...,
	)
	companyRoutes.GET(
		"/:companyId",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getCompany"),
				middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.company.getCompany"),
			},
			defaultMiddlewares,
			module.Binder.BindGetCompany(module.Controller.GetCompany),
		)...,
	)
	companyRoutes.PATCH(
		"",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateCompany"),
				middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.company.updateCompany"),
			},
			defaultMiddlewares,
			module.Binder.BindUpdateCompany(module.Controller.UpdateCompany),
		)...,
	)
	companyRoutes.GET(
		"/:companyId/members",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getCompanyMembers"),
				middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.company.getCompanyMembers"),
			},
			defaultMiddlewares,
			module.Binder.BindGetCompanyMembers(module.Controller.GetCompanyMembers),
		)...,
	)
	companyRoutes.POST(
		"/members",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "addCompanyMember"),
				middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.company.addCompanyMember"),
			},
			defaultMiddlewares,
			module.Binder.BindAddCompanyMember(module.Controller.AddCompanyMember),
		)...,
	)
	companyRoutes.PATCH(
		"/members",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateCompanyMember"),
				middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.company.updateCompanyMember"),
			},
			defaultMiddlewares,
			module.Binder.BindUpdateCompanyMember(module.Controller.UpdateCompanyMember),
		)...,
	)
	companyRoutes.DELETE(
		"/members",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteCompanyMember"),
				middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.company.deleteCompanyMember"),
			},
			defaultMiddlewares,
			module.Binder.BindDeleteCompanyMember(module.Controller.DeleteCompanyMember),
		)...,
	)
}
