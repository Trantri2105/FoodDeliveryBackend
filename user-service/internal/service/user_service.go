package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/pkg/jwt"
)

type UserService interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	Login(ctx context.Context, email, password string) (userId int, accessId string, role string, err error)
	GetUserById(ctx context.Context, userId int) (model.User, error)
	UpdateUserById(ctx context.Context, user model.User) (model.User, error)
	UpdatePasswordById(ctx context.Context, userId int, oldPassword, newPassword string) error
	GetUserList(ctx context.Context, limit, offset int) ([]model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	jwtUtils jwt.Utils
}

func (u *userService) GetUserList(ctx context.Context, limit, offset int) ([]model.User, error) {
	return u.userRepo.GetUserList(ctx, limit, offset)
}

func (u *userService) Register(ctx context.Context, user model.User) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("User service, register user err: %v", err)
		return model.User{}, err
	}
	user.Password = string(hash)
	if user.Role != model.CUSTOMER && user.Role != model.SHIPPER {
		return model.User{}, errors.New("user role must be customer or shipper")
	}
	return u.userRepo.CreateUser(ctx, user)
}

func (u *userService) Login(ctx context.Context, email, password string) (userId int, accessId string, role string, err error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, "", "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, "", "", errors.New("wrong password")
	}
	token, err := u.jwtUtils.CreateToken(user.UserId, user.Role)
	if err != nil {
		return 0, "", "", err
	}
	return user.UserId, token, user.Role, nil
}

func (u *userService) GetUserById(ctx context.Context, userId int) (model.User, error) {
	return u.userRepo.GetUserById(ctx, userId)
}

func (u *userService) UpdateUserById(ctx context.Context, user model.User) (model.User, error) {
	return u.userRepo.UpdateUserById(ctx, user)
}

func (u *userService) UpdatePasswordById(ctx context.Context, userId int, oldPassword, newPassword string) error {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("wrong current password")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	_, err = u.userRepo.UpdateUserById(ctx, user)
	return err
}

func NewUserService(userRepo repository.UserRepository, utils jwt.Utils) UserService {
	return &userService{userRepo: userRepo, jwtUtils: utils}
}
