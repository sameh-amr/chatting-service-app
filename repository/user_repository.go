package repository

import (
    "chatting-service-app/models"
    "chatting-service-app/db"
    "errors"
    "gorm.io/gorm"
)

type UserRepository struct {}

func NewUserRepository() *UserRepository {
    return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *models.User) error {
    result := db.DB.Create(user)
    return result.Error
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    result := db.DB.Where("username = ?", username).First(&user)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &user, result.Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    result := db.DB.Where("email = ?", email).First(&user)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &user, result.Error
}
