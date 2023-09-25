package usecase

import (
	"errors"

	"github.com/benbjohnson/clock"
)

// Repositories - интерфейс инкапсулирующий в себе все репозитории
type Repositories interface {
	LoginRepository
	UserStatesHandlerRepository
	TableRepostiory
	AuthRepository
}

type usecase struct {
	loginRepository          LoginRepository
	registerStatesRepository UserStatesHandlerRepository
	tableRepostiory          TableRepostiory
	authRepository           AuthRepository
	clock                    clock.Clock
}

// NewUsecase - конструктор для usecase
func NewUsecase(repositories Repositories) usecase {
	return usecase{
		loginRepository:          repositories,
		registerStatesRepository: repositories,
		tableRepostiory:          repositories,
		authRepository:           repositories,
		clock:                    clock.New(),
	}
}

const DefaultErrorMessage string = "Упс, что-то пошло не так"

// ErrNotFound - ошибка "не найдено"
var ErrNotFound = errors.New("not found error")

var ErrTeamIdIsUsed = errors.New("uniq constraint")
