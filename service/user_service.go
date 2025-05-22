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

func (s *UserService) Authenticate(email, password string) (*models.User, error) {
    user, err := s.repo.GetUserByEmail(email)
    if err != nil || user == nil {
        return nil, errors.New("invalid email or password")
    }
    if !utils.CheckPasswordHash(password, user.Password) {
        return nil, errors.New("invalid email or password")
    }
    return user, nil
}

func (s *UserService) SignUpAndToken(username, email, password string) (string, error) {
    err := s.SignUp(username, email, password)
    if err != nil {
        return "", err
    }
    user, err := s.Authenticate(email, password)
    if err != nil || user == nil {
        return "", err
    }
    token, err := utils.GenerateJWT(user.ID.String())
    if err != nil {
        return "", err
    }
    return token, nil
}

func (s *UserService) LoginAndToken(email, password string) (string, error) {
    user, err := s.Authenticate(email, password)
    if err != nil || user == nil {
        return "", err
    }
    token, err := utils.GenerateJWT(user.ID.String())
    if err != nil {
        return "", err
    }
    return token, nil
}

func (s *UserService) SetOnlineStatus(userID string, isOnline bool) error {
    return s.repo.SetOnlineStatus(userID, isOnline)
}

func (s *UserService) GetOnlineUsers() ([]models.User, error) {
    return s.repo.GetOnlineUsers()
}
