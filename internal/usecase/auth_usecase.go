package usecase

import (
	"errors"

	"github.com/example/devfolio-api/internal/domain/repositories"
	pkgAuth "github.com/example/devfolio-api/pkg/auth"
)

type AuthUsecase struct {
	userRepo       repositories.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

func NewAuthUsecase(userRepo repositories.UserRepository, jwtSecret string, jwtExpiryHours int) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, jwtSecret: jwtSecret, jwtExpiryHours: jwtExpiryHours}
}

func (u *AuthUsecase) Login(email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}
	if err := pkgAuth.CheckPassword(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid credentials")
	}
	return pkgAuth.GenerateToken(u.jwtSecret, user.ID, user.Email, user.DisplayName, u.jwtExpiryHours)
}
