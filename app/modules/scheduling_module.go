package modules

import (
	ai "github.com/HiIamJeff67/shift-hero-backend/app/ai"
	binders "github.com/HiIamJeff67/shift-hero-backend/app/binders"
	controllers "github.com/HiIamJeff67/shift-hero-backend/app/controllers"
	models "github.com/HiIamJeff67/shift-hero-backend/app/models"
	repositories "github.com/HiIamJeff67/shift-hero-backend/app/models/repositories"
	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
	util "github.com/HiIamJeff67/shift-hero-backend/app/util"
)

type SchedulingModule struct {
	Binder     binders.SchedulingBinderInterface
	Controller controllers.SchedulingControllerInterface
}

func NewSchedulingModule() *SchedulingModule {
	schedulingService := services.NewSchedulingService(models.DB)
	insightGenerator, insightGeneratorError := ai.NewOpenRouterScheduleInsightGenerator(
		util.GetEnv("OPEN_ROUTER_API_KEY", ""),
		util.GetEnv("OPEN_ROUTER_MODEL", "openai/gpt-oss-20b:free"),
	)
	schedulingInsightService := services.NewSchedulingInsightService(
		models.DB,
		repositories.NewUserAccountRepository(),
		insightGenerator,
		insightGeneratorError,
	)
	schedulingBinder := binders.NewSchedulingBinder()
	schedulingController := controllers.NewSchedulingController(
		schedulingService,
		schedulingInsightService,
	)

	return &SchedulingModule{
		Binder:     schedulingBinder,
		Controller: schedulingController,
	}
}
