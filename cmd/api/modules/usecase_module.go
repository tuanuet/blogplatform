package modules

import (
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"go.uber.org/fx"
)

// UseCaseModule provides application use case dependencies
var UseCaseModule = fx.Module("usecase",
	fx.Provide(
		usecase.NewRoleUseCase,
		usecase.NewPermissionUseCase,
		usecase.NewProfileUseCase,
		usecase.NewHealthUseCase,
		usecase.NewCategoryUseCase,
		usecase.NewTagUseCase,
		usecase.NewSubscriptionUseCase,
		usecase.NewCommentUseCase,
		usecase.NewRankingUseCase,
	),
)
