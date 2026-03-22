package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/domain/repositories"
	"github.com/example/devfolio-api/pkg/utils"
)

type CreatePostInput struct {
	Title         string
	Summary       string
	Content       string
	CoverImageURL string
	Status        string
	TagNames      []string
	CreatedBy     uint
}

type UpdatePostInput struct {
	ID            uint
	Title         string
	Summary       string
	Content       string
	CoverImageURL string
	Status        string
	TagNames      []string
	UpdatedBy     uint
}

type PostUsecase struct {
	postRepo repositories.PostRepository
	tagRepo  repositories.TagRepository
}

func NewPostUsecase(postRepo repositories.PostRepository, tagRepo repositories.TagRepository) *PostUsecase {
	return &PostUsecase{postRepo: postRepo, tagRepo: tagRepo}
}

func (u *PostUsecase) ListPublished() ([]entities.Post, error) {
	return u.postRepo.ListPublished()
}

func (u *PostUsecase) GetPublishedBySlug(slug string) (*entities.Post, error) {
	return u.postRepo.GetPublishedBySlug(slug)
}

func (u *PostUsecase) Create(input CreatePostInput) (*entities.Post, error) {
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
		return nil, errors.New("title and content are required")
	}
	status := normalizeStatus(input.Status)
	post := &entities.Post{
		Title:         strings.TrimSpace(input.Title),
		Slug:          utils.Slugify(input.Title),
		Summary:       strings.TrimSpace(input.Summary),
		Content:       input.Content,
		CoverImageURL: strings.TrimSpace(input.CoverImageURL),
		Status:        status,
		CreatedBy:     input.CreatedBy,
		UpdatedBy:     input.CreatedBy,
	}
	if status == "published" {
		now := time.Now()
		post.PublishedAt = &now
	}
	if len(input.TagNames) > 0 {
		tags, err := u.resolveTags(input.TagNames)
		if err != nil {
			return nil, err
		}
		post.Tags = tags
	}
	if err := u.postRepo.Create(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (u *PostUsecase) Update(input UpdatePostInput) (*entities.Post, error) {
	existing, err := u.postRepo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("post not found")
	}
	existing.Title = strings.TrimSpace(input.Title)
	existing.Slug = utils.Slugify(input.Title)
	existing.Summary = strings.TrimSpace(input.Summary)
	existing.Content = input.Content
	existing.CoverImageURL = strings.TrimSpace(input.CoverImageURL)
	existing.Status = normalizeStatus(input.Status)
	existing.UpdatedBy = input.UpdatedBy
	if existing.Status == "published" && existing.PublishedAt == nil {
		now := time.Now()
		existing.PublishedAt = &now
	}
	if len(input.TagNames) > 0 {
		tags, err := u.resolveTags(input.TagNames)
		if err != nil {
			return nil, err
		}
		existing.Tags = tags
	}
	if err := u.postRepo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *PostUsecase) Delete(id uint) error {
	return u.postRepo.Delete(id)
}

func (u *PostUsecase) resolveTags(tagNames []string) ([]entities.Tag, error) {
	resolved := make([]entities.Tag, 0, len(tagNames))
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		tag := entities.Tag{Name: name, Slug: utils.Slugify(name)}
		existing, err := u.tagRepo.GetByNames([]string{name})
		if err != nil {
			return nil, err
		}
		if len(existing) == 0 {
			if err := u.tagRepo.Create(&tag); err != nil {
				return nil, err
			}
			resolved = append(resolved, tag)
			continue
		}
		resolved = append(resolved, existing[0])
	}
	return resolved, nil
}

func normalizeStatus(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "draft", "published", "archived":
		return s
	default:
		return "draft"
	}
}
