package usecase

import (
	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/domain/repositories"
)

type ProjectUsecase struct {
	repo repositories.ProjectRepository
}

func NewProjectUsecase(repo repositories.ProjectRepository) *ProjectUsecase {
	return &ProjectUsecase{repo: repo}
}

func (u *ProjectUsecase) ListFeatured() ([]entities.Project, error) {
	return u.repo.ListFeatured()
}
