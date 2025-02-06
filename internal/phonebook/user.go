package phonebook

import (
	"context"
	"errors"
)

type User struct {
	ID       int
	Email    string
	Password string
}

type UserRepository interface {
	NewUser(context.Context, *User) error
	GetUserByEmail(context.Context, string) (*User, error)
}

type UserService struct {
	repo UserRepository
}

func (s *UserService) Register(ctx context.Context, user *User) error {
	user, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if user != nil {
		return errors.New("email already registered")
	}

	err = s.repo.NewUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, loginUser *User) error {
	user, err := s.repo.GetUserByEmail(ctx, loginUser.Email)
	if err != nil {
		return err
	}

	if user == nil || user.Password != loginUser.Password {
		return errors.New("incorrect email or password")
	}

	return nil

}
