package modules

import (
	binders "github.com/your-org/go-start-monolithic-kit/app/binders"
	controllers "github.com/your-org/go-start-monolithic-kit/app/controllers"
	models "github.com/your-org/go-start-monolithic-kit/app/models"
	repositories "github.com/your-org/go-start-monolithic-kit/app/models/repositories"
	services "github.com/your-org/go-start-monolithic-kit/app/services"
)

type UserSettingModule struct {
	Binder     binders.UserSettingBinderInterface
	Controller controllers.UserSettingControllerInterface
}

func NewUserSettingModule() *UserSettingModule {
	userSettingRepository := repositories.NewUserSettingRepository()

	userSettingService := services.NewUserSettingService(
		models.DB,
		userSettingRepository,
	)

	userSettingBinder := binders.NewUserSettingBinder()

	userSettingController := controllers.NewUserSettingController(
		userSettingService,
	)

	return &UserSettingModule{
		Binder:     userSettingBinder,
		Controller: userSettingController,
	}
}
