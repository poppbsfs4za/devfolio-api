package repository

import (
	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/infrastructure/persistence/gormmodel"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) ListFeatured() ([]entities.Project, error) {
	var models []gormmodel.Project
	if err := r.db.Where("featured = ?", true).Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]entities.Project, 0, len(models))
	for _, p := range models {
		result = append(result, entities.Project{
			ID: p.ID, Name: p.Name, Slug: p.Slug, Description: p.Description,
			RepoURL: p.RepoURL, DemoURL: p.DemoURL, TechStack: p.TechStack,
			Featured: p.Featured, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
		})
	}
	return result, nil
}
