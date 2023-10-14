package usecase

import (
	"errors"

	"github.com/benbjohnson/clock"
)

// Repositories - интерфейс инкапсулирующий в себе все репозитории
type Repositories interface {
	LoginRepository
	UserStatesHandlerRepository
	AuthRepository
	StartRepository
	RegisterRepository
	RegistrationsRepository
	RegisterAvailableRepository
}

type usecase struct {
	loginRepository             LoginRepository
	registerStatesRepository    UserStatesHandlerRepository
	authRepository              AuthRepository
	startRepository             StartRepository
	registerRepository          RegisterRepository
	registrationsRepository     RegistrationsRepository
	registerAvailableRepository RegisterAvailableRepository
	clock                       clock.Clock
}

// NewUsecase - конструктор для usecase
func NewUsecase(repositories Repositories) usecase {
	return usecase{
		loginRepository:             repositories,
		registerStatesRepository:    repositories,
		authRepository:              repositories,
		startRepository:             repositories,
		registerRepository:          repositories,
		registrationsRepository:     repositories,
		registerAvailableRepository: repositories,

		clock: clock.New(),
	}
}

const DefaultErrorMessage string = "Упс, что-то пошло не так"

// ErrNotFound - ошибка "не найдено"
var ErrNotFound = errors.New("not found error")

var ErrTeamIdIsUsed = errors.New("uniq constraint")
