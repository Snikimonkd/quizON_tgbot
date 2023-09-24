package usecase

import (
	"errors"

	"github.com/benbjohnson/clock"
)

// Repositories - интерфейс инкапсулирующий в себе все репозитории
type Repositories interface {
	GamesRepository
	CreateRepository
	LoginRepository
	CheckAuthRepository
	RegisterRepository
	//	ListRepository
	RegisterStatesRepository
}

type usecase struct {
	gamesRepository     GamesRepository
	createRepository    CreateRepository
	loginRepository     LoginRepository
	checkAuthRepository CheckAuthRepository
	registerRepository  RegisterRepository
	//	listRepository           ListRepository
	registerStatesRepository RegisterStatesRepository
	clock                    clock.Clock
}

// NewUsecase - конструктор для usecase
func NewUsecase(repositories Repositories) usecase {
	return usecase{
		gamesRepository:     repositories,
		createRepository:    repositories,
		loginRepository:     repositories,
		checkAuthRepository: repositories,
		registerRepository:  repositories,
		//		listRepository:           repositories,
		registerStatesRepository: repositories,
		clock:                    clock.New(),
	}
}

const DefaultErrorMessage string = "Упс, что-то пошло не так"

// ErrNotFound - ошибка "не найдено"
var ErrNotFound = errors.New("not found error")

var ErrTeamIdIsUsed = errors.New("uniq constraint")
