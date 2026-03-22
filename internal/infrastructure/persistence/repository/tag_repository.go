package repository

import (
	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/infrastructure/persistence/gormmodel"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) List() ([]entities.Tag, error) {
	var models []gormmodel.Tag
	if err := r.db.Order("name asc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]entities.Tag, 0, len(models))
	for _, tag := range models {
		result = append(result, entities.Tag{ID: tag.ID, Name: tag.Name, Slug: tag.Slug, CreatedAt: tag.CreatedAt, UpdatedAt: tag.UpdatedAt})
	}
	return result, nil
}

func (r *TagRepository) GetByNames(names []string) ([]entities.Tag, error) {
	if len(names) == 0 {
		return []entities.Tag{}, nil
	}
	var models []gormmodel.Tag
	if err := r.db.Where("name IN ?", names).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]entities.Tag, 0, len(models))
	for _, tag := range models {
		result = append(result, entities.Tag{ID: tag.ID, Name: tag.Name, Slug: tag.Slug, CreatedAt: tag.CreatedAt, UpdatedAt: tag.UpdatedAt})
	}
	return result, nil
}

func (r *TagRepository) Create(tag *entities.Tag) error {
	model := gormmodel.Tag{Name: tag.Name, Slug: tag.Slug}
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	tag.ID = model.ID
	tag.CreatedAt = model.CreatedAt
	tag.UpdatedAt = model.UpdatedAt
	return nil
}
