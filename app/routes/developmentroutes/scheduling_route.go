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

func configureDevelopmentSchedulingRoutes() {
	module := modules.NewSchedulingModule()

	routes := DevelopmentRouterGroup.Group("/companies")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}

	routes.POST("/shiftRequirements", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createShiftRequirement"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.createShiftRequirement"),
		},
		defaultMiddlewares,
		module.Binder.BindCreateShiftRequirement(module.Controller.CreateShiftRequirement),
	)...)
	routes.GET("/:companyId/shiftRequirements", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getShiftRequirements"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.getShiftRequirements"),
		},
		defaultMiddlewares,
		module.Binder.BindGetShiftRequirements(module.Controller.GetShiftRequirements),
	)...)
	routes.PATCH("/shiftRequirements", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateShiftRequirement"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.updateShiftRequirement"),
		},
		defaultMiddlewares,
		module.Binder.BindUpdateShiftRequirement(module.Controller.UpdateShiftRequirement),
	)...)
	routes.DELETE("/shiftRequirements", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteShiftRequirement"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.deleteShiftRequirement"),
		},
		defaultMiddlewares,
		module.Binder.BindDeleteShiftRequirement(module.Controller.DeleteShiftRequirement),
	)...)

	routes.PUT("/availabilitySlots", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "upsertAvailabilitySlots"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.upsertAvailabilitySlots"),
		},
		defaultMiddlewares,
		module.Binder.BindUpsertAvailabilitySlots(module.Controller.UpsertAvailabilitySlots),
	)...)
	routes.GET("/:companyId/availabilitySlots", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAvailabilitySlots"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.getAvailabilitySlots"),
		},
		defaultMiddlewares,
		module.Binder.BindGetAvailabilitySlots(module.Controller.GetAvailabilitySlots),
	)...)
	routes.DELETE("/availabilitySlots", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteAvailabilitySlot"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.deleteAvailabilitySlot"),
		},
		defaultMiddlewares,
		module.Binder.BindDeleteAvailabilitySlot(module.Controller.DeleteAvailabilitySlot),
	)...)

	routes.POST("/assignments/generate", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "generateAssignments"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.generateAssignments"),
		},
		defaultMiddlewares,
		module.Binder.BindGenerateAssignments(module.Controller.GenerateAssignments),
	)...)
	routes.PUT("/assignments", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "replaceAssignments"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.replaceAssignments"),
		},
		defaultMiddlewares,
		module.Binder.BindReplaceAssignments(module.Controller.ReplaceAssignments),
	)...)
	routes.POST("/assignments/claim", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "claimAssignment"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.claimAssignment"),
		},
		defaultMiddlewares,
		module.Binder.BindClaimAssignment(module.Controller.ClaimAssignment),
	)...)
	routes.GET("/:companyId/assignments", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAssignments"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.getAssignments"),
		},
		defaultMiddlewares,
		module.Binder.BindGetAssignments(module.Controller.GetAssignments),
	)...)

	routes.POST("/swapRequests", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createSwapRequest"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.createSwapRequest"),
		},
		defaultMiddlewares,
		module.Binder.BindCreateSwapRequest(module.Controller.CreateSwapRequest),
	)...)
	routes.GET("/:companyId/swapRequests", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getSwapRequests"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.getSwapRequests"),
		},
		defaultMiddlewares,
		module.Binder.BindGetSwapRequests(module.Controller.GetSwapRequests),
	)...)
	routes.POST("/swapRequests/claim", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "claimSwapRequest"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.claimSwapRequest"),
		},
		defaultMiddlewares,
		module.Binder.BindClaimSwapRequest(module.Controller.ClaimSwapRequest),
	)...)
	routes.POST("/swapRequests/approve", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "approveSwapRequest"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.approveSwapRequest"),
		},
		defaultMiddlewares,
		module.Binder.BindApproveSwapRequest(module.Controller.ApproveSwapRequest),
	)...)
	routes.POST("/swapRequests/cancel", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "cancelSwapRequest"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.cancelSwapRequest"),
		},
		defaultMiddlewares,
		module.Binder.BindCancelSwapRequest(module.Controller.CancelSwapRequest),
	)...)

	routes.GET("/:companyId/schedulePublications", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getSchedulePublication"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.getSchedulePublication"),
		},
		defaultMiddlewares,
		module.Binder.BindGetSchedulePublication(module.Controller.GetSchedulePublication),
	)...)
	routes.PUT("/schedulePublications", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "upsertSchedulePublication"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.upsertSchedulePublication"),
		},
		defaultMiddlewares,
		module.Binder.BindUpsertSchedulePublication(module.Controller.UpsertSchedulePublication),
	)...)

	routes.GET("/:companyId/scheduleSettings", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getScheduleSettings"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.getScheduleSettings"),
		},
		defaultMiddlewares,
		module.Binder.BindGetCompanySettings(module.Controller.GetCompanySettings),
	)...)
	routes.PATCH("/scheduleSettings", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateScheduleSettings"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.updateScheduleSettings"),
		},
		defaultMiddlewares,
		module.Binder.BindUpdateCompanySettings(module.Controller.UpdateCompanySettings),
	)...)

	routes.GET("/:companyId/settings", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getCompanySettings"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.getCompanySettings"),
		},
		defaultMiddlewares,
		module.Binder.BindGetCompanySettings(module.Controller.GetCompanySettings),
	)...)
	routes.PATCH("/settings", middlewares.RepositionMiddleware(
		[]gin.HandlerFunc{
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateCompanySettings"),
			middlewares.ApplyMeterMiddleware(otel.Meter(constants.ServiceName), "server.requests.scheduling.updateCompanySettings"),
		},
		defaultMiddlewares,
		module.Binder.BindUpdateCompanySettings(module.Controller.UpdateCompanySettings),
	)...)
}
