package modules

import (
	binders "github.com/your-org/go-start-monolithic-kit/app/binders"
	controllers "github.com/your-org/go-start-monolithic-kit/app/controllers"
	models "github.com/your-org/go-start-monolithic-kit/app/models"
	repositories "github.com/your-org/go-start-monolithic-kit/app/models/repositories"
	services "github.com/your-org/go-start-monolithic-kit/app/services"
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
