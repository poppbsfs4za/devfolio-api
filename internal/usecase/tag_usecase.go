package usecase

import (
	"strings"

	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/domain/repositories"
	"github.com/example/devfolio-api/pkg/utils"
)

type TagUsecase struct {
	repo repositories.TagRepository
}

func NewTagUsecase(repo repositories.TagRepository) *TagUsecase {
	return &TagUsecase{repo: repo}
}

func (u *TagUsecase) List() ([]entities.Tag, error) {
	return u.repo.List()
}

func (u *TagUsecase) Create(name string) (*entities.Tag, error) {
	tag := &entities.Tag{Name: strings.TrimSpace(name), Slug: utils.Slugify(name)}
	if err := u.repo.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}
