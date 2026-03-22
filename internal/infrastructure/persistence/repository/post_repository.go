package repository

import (
	"errors"

	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/infrastructure/persistence/gormmodel"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) ListPublished() ([]entities.Post, error) {
	var models []gormmodel.Post
	if err := r.db.Preload("Tags").Where("status = ?", "published").Order("published_at desc nulls last, created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toPostEntities(models), nil
}

func (r *PostRepository) GetPublishedBySlug(slug string) (*entities.Post, error) {
	var model gormmodel.Post
	if err := r.db.Preload("Tags").Where("slug = ? AND status = ?", slug, "published").First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	entity := toPostEntity(model)
	return &entity, nil
}

func (r *PostRepository) GetByID(id uint) (*entities.Post, error) {
	var model gormmodel.Post
	if err := r.db.Preload("Tags").First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	entity := toPostEntity(model)
	return &entity, nil
}

func (r *PostRepository) Create(post *entities.Post) error {
	model := gormmodel.Post{
		Title:         post.Title,
		Slug:          post.Slug,
		Summary:       post.Summary,
		Content:       post.Content,
		CoverImageURL: post.CoverImageURL,
		Status:        post.Status,
		PublishedAt:   post.PublishedAt,
		CreatedBy:     post.CreatedBy,
		UpdatedBy:     post.UpdatedBy,
	}
	if len(post.Tags) > 0 {
		model.Tags = toTagModels(post.Tags)
	}
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	post.ID = model.ID
	post.CreatedAt = model.CreatedAt
	post.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *PostRepository) Update(post *entities.Post) error {
	var model gormmodel.Post
	if err := r.db.Preload("Tags").First(&model, post.ID).Error; err != nil {
		return err
	}
	model.Title = post.Title
	model.Slug = post.Slug
	model.Summary = post.Summary
	model.Content = post.Content
	model.CoverImageURL = post.CoverImageURL
	model.Status = post.Status
	model.PublishedAt = post.PublishedAt
	model.UpdatedBy = post.UpdatedBy
	if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&model).Association("Tags").Replace(toTagModels(post.Tags)); err != nil {
		return err
	}
	return r.db.Save(&model).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&gormmodel.Post{}, id).Error
}

func toPostEntities(models []gormmodel.Post) []entities.Post {
	result := make([]entities.Post, 0, len(models))
	for _, model := range models {
		result = append(result, toPostEntity(model))
	}
	return result
}

func toPostEntity(model gormmodel.Post) entities.Post {
	tags := make([]entities.Tag, 0, len(model.Tags))
	for _, tag := range model.Tags {
		tags = append(tags, entities.Tag{ID: tag.ID, Name: tag.Name, Slug: tag.Slug, CreatedAt: tag.CreatedAt, UpdatedAt: tag.UpdatedAt})
	}
	return entities.Post{
		ID:            model.ID,
		Title:         model.Title,
		Slug:          model.Slug,
		Summary:       model.Summary,
		Content:       model.Content,
		CoverImageURL: model.CoverImageURL,
		Status:        model.Status,
		PublishedAt:   model.PublishedAt,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		CreatedBy:     model.CreatedBy,
		UpdatedBy:     model.UpdatedBy,
		Tags:          tags,
	}
}

func toTagModels(tags []entities.Tag) []gormmodel.Tag {
	result := make([]gormmodel.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, gormmodel.Tag{BaseModel: gormmodel.BaseModel{ID: tag.ID}, Name: tag.Name, Slug: tag.Slug})
	}
	return result
}
