package modules

import (
	binders "github.com/HiIamJeff67/shift-hero-backend/app/binders"
	controllers "github.com/HiIamJeff67/shift-hero-backend/app/controllers"
	models "github.com/HiIamJeff67/shift-hero-backend/app/models"
	repositories "github.com/HiIamJeff67/shift-hero-backend/app/models/repositories"
	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
)

type CompanyModule struct {
	Binder     binders.CompanyBinderInterface
	Controller controllers.CompanyControllerInterface
}

func NewCompanyModule() *CompanyModule {
	companyRepository := repositories.NewCompanyRepository()
	usersToCompaniesRepository := repositories.NewUsersToCompaniesRepository()
	userRepository := repositories.NewUserRepository()

	companyService := services.NewCompanyService(
		models.DB,
		companyRepository,
		usersToCompaniesRepository,
		userRepository,
	)

	companyBinder := binders.NewCompanyBinder()

	companyController := controllers.NewCompanyController(
		companyService,
	)

	return &CompanyModule{
		Binder:     companyBinder,
		Controller: companyController,
	}
}
