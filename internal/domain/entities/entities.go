package entities

import "time"

type User struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	DisplayName  string    `json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Profile struct {
	ID          uint      `json:"id"`
	FullName    string    `json:"full_name"`
	Headline    string    `json:"headline"`
	Bio         string    `json:"bio"`
	GitHubURL   string    `json:"github_url"`
	LinkedInURL string    `json:"linkedin_url"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Project struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	RepoURL     string    `json:"repo_url"`
	DemoURL     string    `json:"demo_url"`
	TechStack   string    `json:"tech_stack"`
	Featured    bool      `json:"featured"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Tag struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Summary       string     `json:"summary"`
	Content       string     `json:"content"`
	CoverImageURL string     `json:"cover_image_url"`
	Status        string     `json:"status"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uint       `json:"created_by"`
	UpdatedBy     uint       `json:"updated_by"`
	Tags          []Tag      `json:"tags"`
}
