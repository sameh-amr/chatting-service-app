package service

import (
    "errors"
    "chatting-service-app/models"
    "chatting-service-app/repository"
    "chatting-service-app/utils"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) SignUp(username, email, password string) error {
    if username == "" || email == "" || password == "" {
        return errors.New("all fields are required")
    }

    existingUser, _ := s.repo.GetUserByUsername(username)
    if existingUser != nil {
        return errors.New("username already taken")
    }

    existingEmail, _ := s.repo.GetUserByEmail(email)
    if existingEmail != nil {
        return errors.New("email already registered")
    }

    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        return err
    }

    user := &models.User{
        Username: username,
        Email: email,
        Password: hashedPassword,
    }

    return s.repo.CreateUser(user)
}
