package Usecases

import (
	"context"
	"errors"
	"time"

	"task_manager/Domain"
	"task_manager/Repositories"
	"task_manager/Infrastructure"
)

type UserUsecase interface {
	Register(username, password string) (Domain.User, error)
	Login(username, password string) (Domain.User, error)
	Promote(username string) error
}

type userUsecase struct {
	repo Repositories.UserRepository
	timeout time.Duration
}

func NewUserUsecase(r Repositories.UserRepository) UserUsecase {
	return &userUsecase{repo: r, timeout: 5 * time.Second}
}

func (u *userUsecase) Register(username, password string) (Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	// if first user, make admin
	count, err := u.repo.CountUsers(ctx)
	if err != nil {
		return Domain.User{}, err
	}
	role := "user"
	if count == 0 {
		role = "admin"
	}
	hashed, err := Infrastructure.HashPassword(password)
	if err != nil {
		return Domain.User{}, err
	}
	user := Domain.User{Username: username, Password: hashed, Role: role}
	created, err := u.repo.Create(ctx, user)
	if err != nil {
		return Domain.User{}, err
	}
	created.Password = ""
	return created, nil
}

func (u *userUsecase) Login(username, password string) (Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	user, err := u.repo.FindByUsername(ctx, username)
	if err != nil {
		return Domain.User{}, errors.New("invalid credentials")
	}
	// compare
	if err := Infrastructure.ComparePassword(user.Password, password); err != nil {
		return Domain.User{}, errors.New("invalid credentials")
	}
	user.Password = ""
	return user, nil
}

func (u *userUsecase) Promote(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()
	return u.repo.UpdateRole(ctx, username, "admin")
}
