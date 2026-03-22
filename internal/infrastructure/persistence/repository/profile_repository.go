package repository

import (
	"errors"

	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/infrastructure/persistence/gormmodel"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Get() (*entities.Profile, error) {
	var profile gormmodel.Profile
	if err := r.db.First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entities.Profile{
		ID:          profile.ID,
		FullName:    profile.FullName,
		Headline:    profile.Headline,
		Bio:         profile.Bio,
		GitHubURL:   profile.GitHubURL,
		LinkedInURL: profile.LinkedInURL,
		AvatarURL:   profile.AvatarURL,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}, nil
}

func (r *ProfileRepository) Upsert(profile *entities.Profile) error {
	var existing gormmodel.Profile
	result := r.db.First(&existing)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			model := gormmodel.Profile{
				FullName:    profile.FullName,
				Headline:    profile.Headline,
				Bio:         profile.Bio,
				GitHubURL:   profile.GitHubURL,
				LinkedInURL: profile.LinkedInURL,
				AvatarURL:   profile.AvatarURL,
			}
			return r.db.Create(&model).Error
		}
		return result.Error
	}

	existing.FullName = profile.FullName
	existing.Headline = profile.Headline
	existing.Bio = profile.Bio
	existing.GitHubURL = profile.GitHubURL
	existing.LinkedInURL = profile.LinkedInURL
	existing.AvatarURL = profile.AvatarURL
	return r.db.Save(&existing).Error
}
