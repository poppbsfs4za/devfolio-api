package gormmodel

import "time"

type BaseModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	BaseModel
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	DisplayName  string `gorm:"size:255;not null"`
}

type Profile struct {
	BaseModel
	FullName    string `gorm:"size:255;not null"`
	Headline    string `gorm:"size:255"`
	Bio         string `gorm:"type:text"`
	GitHubURL   string `gorm:"size:500"`
	LinkedInURL string `gorm:"size:500"`
	AvatarURL   string `gorm:"size:500"`
}

type Project struct {
	BaseModel
	Name        string `gorm:"size:255;not null"`
	Slug        string `gorm:"size:255;uniqueIndex;not null"`
	Description string `gorm:"type:text"`
	RepoURL     string `gorm:"size:500"`
	DemoURL     string `gorm:"size:500"`
	TechStack   string `gorm:"type:text"`
	Featured    bool   `gorm:"not null;default:false"`
}

type Tag struct {
	BaseModel
	Name string `gorm:"size:100;uniqueIndex;not null"`
	Slug string `gorm:"size:120;uniqueIndex;not null"`
}

type Post struct {
	BaseModel
	Title         string `gorm:"size:255;not null"`
	Slug          string `gorm:"size:255;uniqueIndex;not null"`
	Summary       string `gorm:"type:text"`
	Content       string `gorm:"type:text;not null"`
	CoverImageURL string `gorm:"size:500"`
	Status        string `gorm:"size:50;index;not null"`
	PublishedAt   *time.Time
	CreatedBy     uint  `gorm:"not null"`
	UpdatedBy     uint  `gorm:"not null"`
	Tags          []Tag `gorm:"many2many:post_tags;"`
}

type PostTag struct {
	PostID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}
