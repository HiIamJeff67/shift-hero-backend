package modules

import (
	binders "github.com/your-org/go-start-monolithic-kit/app/binders"
	configs "github.com/your-org/go-start-monolithic-kit/app/configs"
	controllers "github.com/your-org/go-start-monolithic-kit/app/controllers"
	models "github.com/your-org/go-start-monolithic-kit/app/models"
	repositories "github.com/your-org/go-start-monolithic-kit/app/models/repositories"
	services "github.com/your-org/go-start-monolithic-kit/app/services"
)

type AuthModule struct {
	Binder     binders.AuthBinderInterface
	Controller controllers.AuthControllerInterface
}

func NewAuthModule() *AuthModule {
	userRepository := repositories.NewUserRepository()
	userInfoRepository := repositories.NewUserInfoRepository()
	userAccountRepository := repositories.NewUserAccountRepository()
	userSettingRepository := repositories.NewUserSettingRepository()
	oauthService := services.NewOAuthService(configs.OAuthGoogleConfig)

	authService := services.NewAuthService(
		models.DB,
		userRepository,
		userInfoRepository,
		userAccountRepository,
		userSettingRepository,
		oauthService,
	)

	authBinder := binders.NewAuthBinder()

	authController := controllers.NewAuthController(
		authService,
	)

	return &AuthModule{
		Binder:     authBinder,
		Controller: authController,
	}
}
