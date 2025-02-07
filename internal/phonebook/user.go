package phonebook

import (
	"context"
	"template/internal/common"
)

type User struct {
	ID       int
	Email    string
	Password string
}

type UserRepository interface {
	NewUser(context.Context, *User) (int, error)
	GetUserByEmail(context.Context, string) (*User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(ctx context.Context, user *User) (string, error) {
	currentUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}

	if currentUser != nil {
		return "", common.InvariantError{Message: "email already registered"}
	}

	user.Password, err = common.BcryptHash(user.Password)
	if err != nil {
		return "", err
	}

	id, err := s.repo.NewUser(ctx, user)
	if err != nil {
		return "", err
	}

	token, err := common.JwtGenerate(id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) Login(ctx context.Context, loginUser *User) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, loginUser.Email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", common.InvariantError{Message: "incorrect email or password"}
	}

	ok, err := common.BcryptCompare(user.Password, loginUser.Password)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", common.InvariantError{Message: "incorrect email or password"}
	}

	token, err := common.JwtGenerate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil

}
