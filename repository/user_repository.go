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

func (r *UserRepository) GetAllUsersExcept(exceptID string) ([]models.User, error) {
    var users []models.User
    err := db.DB.Where("id != ?", exceptID).Find(&users).Error
    return users, err
}

func (r *UserRepository) SetOnlineStatus(userID string, isOnline bool) error {
    return db.DB.Model(&models.User{}).
        Where("id = ?", userID).
        Update("is_online", isOnline).Error
}

func (r *UserRepository) GetOnlineUsers() ([]models.User, error) {
    var users []models.User
    err := db.DB.Where("is_online = ?", true).Find(&users).Error
    return users, err
}
