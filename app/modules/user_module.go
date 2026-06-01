package modules

import (
	binders "github.com/HiIamJeff67/shift-hero-backend/app/binders"
	controllers "github.com/HiIamJeff67/shift-hero-backend/app/controllers"
	models "github.com/HiIamJeff67/shift-hero-backend/app/models"
	repositories "github.com/HiIamJeff67/shift-hero-backend/app/models/repositories"
	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
)

type UserModule struct {
	Binder     binders.UserBinderInterface
	Controller controllers.UserControllerInterface
}

func NewUserModule() *UserModule {
	userRepository := repositories.NewUserRepository()

	userService := services.NewUserService(
		models.DB,
		userRepository,
	)

	userBinder := binders.NewUserBinder()

	userController := controllers.NewUserController(
		userService,
	)

	return &UserModule{
		Binder:     userBinder,
		Controller: userController,
	}
}
