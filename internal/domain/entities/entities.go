package entities

import "time"

type User struct {
	ID           uint
	Email        string
	PasswordHash string
	DisplayName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Profile struct {
	ID          uint
	FullName    string
	Headline    string
	Bio         string
	GitHubURL   string
	LinkedInURL string
	AvatarURL   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Project struct {
	ID          uint
	Name        string
	Slug        string
	Description string
	RepoURL     string
	DemoURL     string
	TechStack   string
	Featured    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Tag struct {
	ID        uint
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID            uint
	Title         string
	Slug          string
	Summary       string
	Content       string
	CoverImageURL string
	Status        string
	PublishedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CreatedBy     uint
	UpdatedBy     uint
	Tags          []Tag
}
