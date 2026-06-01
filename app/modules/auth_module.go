package modules

import (
	binders "github.com/HiIamJeff67/shift-hero-backend/app/binders"
	configs "github.com/HiIamJeff67/shift-hero-backend/app/configs"
	controllers "github.com/HiIamJeff67/shift-hero-backend/app/controllers"
	models "github.com/HiIamJeff67/shift-hero-backend/app/models"
	repositories "github.com/HiIamJeff67/shift-hero-backend/app/models/repositories"
	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
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
