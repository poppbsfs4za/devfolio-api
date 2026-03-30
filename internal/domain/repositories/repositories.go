package repositories

import "github.com/example/devfolio-api/internal/domain/entities"

type UserRepository interface {
	GetByEmail(email string) (*entities.User, error)
	Create(user *entities.User) error
}

type ProfileRepository interface {
	Get() (*entities.Profile, error)
	Upsert(profile *entities.Profile) error
}

type ProjectRepository interface {
	ListFeatured() ([]entities.Project, error)
}

type TagRepository interface {
	List() ([]entities.Tag, error)
	GetByNames(names []string) ([]entities.Tag, error)
	Create(tag *entities.Tag) error
}

type PostRepository interface {
	ListPublished() ([]entities.Post, error)
	GetPublishedBySlug(slug string) (*entities.Post, error)

	AdminList() ([]entities.Post, error)
	GetByID(id uint) (*entities.Post, error)

	Create(post *entities.Post) error
	Update(post *entities.Post) error
	Delete(id uint) error
}
