package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
	"week_06_web/app_err"
	"week_06_web/model"
	"week_06_web/repository"
)

type UserRepository interface {
	Create(context.Context, model.User) error
	FindById(context.Context, string) (model.User, error)
	FindByUsername(context.Context, string) (model.User, error)
	UpdatePassword(context.Context, string, string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, username string, password string) (model.User, error) {
	username = strings.TrimSpace(username)
	if _, err := s.repo.FindByUsername(ctx, username); err == nil {
		return model.User{},
			app_err.Conflict("USER_ALREADY_EXISTS", "user already exists", err)
	} else if !errors.Is(err, repository.ErrNotFound) {
		return model.User{}, app_err.Internal(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, app_err.Internal(err)
	}

	user := model.User{
		Id:           strconv.FormatInt(time.Now().UnixNano(), 10),
		Username:     username,
		PasswordHash: string(hash),
		Create_at:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return model.User{}, app_err.Internal(err)
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (model.User, error) {
	user, err := s.repo.FindById(ctx, id)
	if err == nil {
		return user, nil
	} else if errors.Is(err, repository.ErrNotFound) {
		return model.User{}, app_err.NotFound("User not found", err)
	}
	return model.User{}, app_err.Internal(err)
}

func (s *UserService) Login(ctx context.Context, username string, password string) (model.User, error) {
	user, err := s.repo.FindByUsername(ctx, strings.TrimSpace(username))
	if errors.Is(err, repository.ErrNotFound) {
		return model.User{}, app_err.Unauthorized("INVALID_CREDENTIALS", "invalid credentials", nil)
	}
	if err != nil {
		return model.User{}, app_err.Internal(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return model.User{}, app_err.Unauthorized("INVALID_CREDENTIALS", "invalid credentials", nil)
	}

	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id string, oldPassword string, newPassword string) error {
	user, err := s.repo.FindById(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return app_err.NotFound("User not found", err)
	}
	if err != nil {
		return app_err.Internal(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return app_err.Unauthorized("INVALID_CREDENTIALS", "old password is incorrect", nil)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return app_err.Internal(err)
	}

	user.PasswordHash = string(hash)
	if err := s.repo.UpdatePassword(ctx, id, string(hash)); err != nil {
		return app_err.Internal(err)
	}

	return nil
}
