package modules

import (
	binders "github.com/your-org/go-start-monolithic-kit/app/binders"
	"github.com/your-org/go-start-monolithic-kit/app/configs"
	controllers "github.com/your-org/go-start-monolithic-kit/app/controllers"
	models "github.com/your-org/go-start-monolithic-kit/app/models"
	repositories "github.com/your-org/go-start-monolithic-kit/app/models/repositories"
	services "github.com/your-org/go-start-monolithic-kit/app/services"
)

type UserAccountModule struct {
	Binder     binders.UserAccountBinderInterface
	Controller controllers.UserAccountControllerInterface
}

func NewUserAccountModule() *UserAccountModule {
	userRepository := repositories.NewUserRepository()
	userAccountRepository := repositories.NewUserAccountRepository()
	oauthService := services.NewOAuthService(configs.OAuthGoogleConfig)

	userAccountService := services.NewUserAccountService(
		models.DB,
		userRepository,
		userAccountRepository,
		oauthService,
	)

	userAccountBinder := binders.NewUserAccountBinder()

	userAccountController := controllers.NewUserAccountController(
		userAccountService,
	)

	return &UserAccountModule{
		Binder:     userAccountBinder,
		Controller: userAccountController,
	}
}
