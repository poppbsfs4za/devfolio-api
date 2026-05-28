package usecase_test

import (
	"errors"
	"testing"

	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/usecase"
)

// --- mock TagRepository ---

type mockTagRepo struct {
	bySlugs  map[string]entities.Tag
	created  []entities.Tag
	createFn func(tag *entities.Tag) error // optional override
}

func (m *mockTagRepo) List() ([]entities.Tag, error) { return nil, nil }

func (m *mockTagRepo) GetByNames(names []string) ([]entities.Tag, error) { return nil, nil }

func (m *mockTagRepo) GetBySlugs(slugs []string) ([]entities.Tag, error) {
	var result []entities.Tag
	for _, s := range slugs {
		if t, ok := m.bySlugs[s]; ok {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *mockTagRepo) Create(tag *entities.Tag) error {
	if m.createFn != nil {
		return m.createFn(tag)
	}
	tag.ID = uint(len(m.created) + 1)
	m.created = append(m.created, *tag)
	m.bySlugs[tag.Slug] = *tag
	return nil
}

// --- mock PostRepository (no-op) ---

type mockPostRepo struct{}

func (m *mockPostRepo) ListPublished() ([]entities.Post, error)               { return nil, nil }
func (m *mockPostRepo) GetPublishedBySlug(slug string) (*entities.Post, error) { return nil, nil }
func (m *mockPostRepo) AdminList() ([]entities.Post, error)                   { return nil, nil }
func (m *mockPostRepo) GetByID(id uint) (*entities.Post, error)               { return nil, nil }
func (m *mockPostRepo) Create(post *entities.Post) error                      { return nil }
func (m *mockPostRepo) Update(post *entities.Post) error                      { return nil }
func (m *mockPostRepo) Delete(id uint) error                                  { return nil }

// --- helper ---

func newUsecase(tagRepo *mockTagRepo) *usecase.PostUsecase {
	return usecase.NewPostUsecase(&mockPostRepo{}, tagRepo)
}

// --- tests ---

func TestCreate_ReuseExistingTagBySlug(t *testing.T) {
	existing := entities.Tag{ID: 7, Name: "Go", Slug: "go"}
	tagRepo := &mockTagRepo{bySlugs: map[string]entities.Tag{"go": existing}}

	post, err := newUsecase(tagRepo).Create(usecase.CreatePostInput{
		Title:    "Test Post",
		Content:  "body",
		TagNames: []string{"Go"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(post.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(post.Tags))
	}
	if post.Tags[0].ID != 7 {
		t.Errorf("expected reused tag ID 7, got %d", post.Tags[0].ID)
	}
	if len(tagRepo.created) != 0 {
		t.Errorf("expected no new tags created, got %d", len(tagRepo.created))
	}
}

func TestCreate_CreatesNewTagWhenNotFound(t *testing.T) {
	tagRepo := &mockTagRepo{bySlugs: map[string]entities.Tag{}}

	post, err := newUsecase(tagRepo).Create(usecase.CreatePostInput{
		Title:    "Test Post",
		Content:  "body",
		TagNames: []string{"Golang"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(post.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(post.Tags))
	}
	if post.Tags[0].Slug != "golang" {
		t.Errorf("expected slug 'golang', got %q", post.Tags[0].Slug)
	}
	if len(tagRepo.created) != 1 {
		t.Errorf("expected 1 new tag, got %d", len(tagRepo.created))
	}
}

func TestCreate_MixedExistingAndNewTags(t *testing.T) {
	existing := entities.Tag{ID: 3, Name: "Go", Slug: "go"}
	tagRepo := &mockTagRepo{bySlugs: map[string]entities.Tag{"go": existing}}

	post, err := newUsecase(tagRepo).Create(usecase.CreatePostInput{
		Title:    "Test Post",
		Content:  "body",
		TagNames: []string{"Go", "Docker"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(post.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(post.Tags))
	}
	if len(tagRepo.created) != 1 {
		t.Errorf("expected 1 new tag created, got %d", len(tagRepo.created))
	}
	if tagRepo.created[0].Slug != "docker" {
		t.Errorf("expected new tag slug 'docker', got %q", tagRepo.created[0].Slug)
	}
}

// Tag name differs in casing but produces the same slug — must reuse, not re-insert.
func TestCreate_SameSlugDifferentCaseIsReused(t *testing.T) {
	existing := entities.Tag{ID: 5, Name: "ai-agent", Slug: "ai-agent"}
	tagRepo := &mockTagRepo{bySlugs: map[string]entities.Tag{"ai-agent": existing}}

	post, err := newUsecase(tagRepo).Create(usecase.CreatePostInput{
		Title:    "Test Post",
		Content:  "body",
		TagNames: []string{"AI-Agent"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.Tags[0].ID != 5 {
		t.Errorf("expected reused tag ID 5, got %d", post.Tags[0].ID)
	}
	if len(tagRepo.created) != 0 {
		t.Errorf("expected no new tags created, got %d", len(tagRepo.created))
	}
}

// Simulates a concurrent insert winning the race: Create returns a duplicate-key error,
// resolveTags must fall back to GetBySlugs and return the winning row.
func TestCreate_RaceConditionFallback(t *testing.T) {
	existing := entities.Tag{ID: 9, Name: "rust", Slug: "rust"}
	tagRepo := &mockTagRepo{bySlugs: map[string]entities.Tag{}}

	tagRepo.createFn = func(tag *entities.Tag) error {
		tagRepo.bySlugs["rust"] = existing
		return errors.New(`ERROR: duplicate key value violates unique constraint "idx_tags_slug"`)
	}

	post, err := newUsecase(tagRepo).Create(usecase.CreatePostInput{
		Title:    "Test Post",
		Content:  "body",
		TagNames: []string{"rust"},
	})
	if err != nil {
		t.Fatalf("unexpected error after race fallback: %v", err)
	}
	if post.Tags[0].ID != 9 {
		t.Errorf("expected fallback tag ID 9, got %d", post.Tags[0].ID)
	}
}

func TestCreate_EmptyTagNames(t *testing.T) {
	tagRepo := &mockTagRepo{bySlugs: map[string]entities.Tag{}}

	post, err := newUsecase(tagRepo).Create(usecase.CreatePostInput{
		Title:    "Test Post",
		Content:  "body",
		TagNames: []string{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(post.Tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(post.Tags))
	}
}
