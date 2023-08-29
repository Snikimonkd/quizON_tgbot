package usecase

import "errors"

// Repositories - интерфейс инкапсулирующий в себе все репозитории
type Repositories interface {
	GamesRepository
	CreateRepository
	LoginRepository
	CheckAuthRepository
	RegisterRepository
	ListRepository
}

type usecase struct {
	gamesRepository     GamesRepository
	createRepository    CreateRepository
	loginRepository     LoginRepository
	checkAuthRepository CheckAuthRepository
	registerRepository  RegisterRepository
	listRepository      ListRepository
}

// NewUsecase - конструктор для usecase
func NewUsecase(repositories Repositories) usecase {
	return usecase{
		gamesRepository:     repositories,
		createRepository:    repositories,
		loginRepository:     repositories,
		checkAuthRepository: repositories,
		registerRepository:  repositories,
		listRepository:      repositories,
	}
}

// ErrNotFound - ошибка "не найдено"
var ErrNotFound = errors.New("not found error")
