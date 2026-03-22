package usecase

import (
	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/domain/repositories"
)

type ProfileUsecase struct {
	repo repositories.ProfileRepository
}

func NewProfileUsecase(repo repositories.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{repo: repo}
}

func (u *ProfileUsecase) Get() (*entities.Profile, error) {
	return u.repo.Get()
}

func (u *ProfileUsecase) Upsert(profile *entities.Profile) error {
	return u.repo.Upsert(profile)
}
