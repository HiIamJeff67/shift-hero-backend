package modules

import (
	binders "github.com/HiIamJeff67/shift-hero-backend/app/binders"
	controllers "github.com/HiIamJeff67/shift-hero-backend/app/controllers"
	models "github.com/HiIamJeff67/shift-hero-backend/app/models"
	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
)

type SchedulingModule struct {
	Binder     binders.SchedulingBinderInterface
	Controller controllers.SchedulingControllerInterface
}

func NewSchedulingModule() *SchedulingModule {
	schedulingService := services.NewSchedulingService(models.DB)
	schedulingBinder := binders.NewSchedulingBinder()
	schedulingController := controllers.NewSchedulingController(schedulingService)

	return &SchedulingModule{
		Binder:     schedulingBinder,
		Controller: schedulingController,
	}
}
