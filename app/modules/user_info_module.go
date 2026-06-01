package modules

import (
	binders "github.com/your-org/go-start-monolithic-kit/app/binders"
	controllers "github.com/your-org/go-start-monolithic-kit/app/controllers"
	models "github.com/your-org/go-start-monolithic-kit/app/models"
	repositories "github.com/your-org/go-start-monolithic-kit/app/models/repositories"
	services "github.com/your-org/go-start-monolithic-kit/app/services"
)

type UserInfoModule struct {
	Binder     binders.UserInfoBinderInterface
	Controller controllers.UserInfoControllerInterface
}

func NewUserInfoModule() *UserInfoModule {
	userInfoRepository := repositories.NewUserInfoRepository()

	userInfoService := services.NewUserInfoService(
		models.DB,
		userInfoRepository,
	)

	userInfoBinder := binders.NewUserInfoBinder()

	userInfoController := controllers.NewUserInfoController(
		userInfoService,
	)

	return &UserInfoModule{
		Binder:     userInfoBinder,
		Controller: userInfoController,
	}
}
