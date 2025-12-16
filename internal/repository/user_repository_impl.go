package repository

import (
	"Gobank/internal/domain"
	"Gobank/internal/repository/models"
	"context"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User, hashedPassword string) (*domain.User, error) {
	userModel := &models.UserModel{
		Username:       user.Username,
		FullName:       user.FullName,
		Email:          user.Email,
		HashedPassword: hashedPassword,
	}
	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return nil, err
	}
	user.ID = userModel.ID
	user.CreatedAt = userModel.CreatedAt
	return user, nil
}

func (r *userRepositoryImpl) GetById(ctx context.Context, id int64) (*domain.User, error) {
	var userModel models.UserModel
	if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        userModel.ID,
		Username:  userModel.Username,
		FullName:  userModel.FullName,
		Email:     userModel.Email,
		CreatedAt: userModel.CreatedAt,
	}, nil
}

func (r *userRepositoryImpl) GetByUserName(ctx context.Context, username string) (*domain.User, string, error) {
	var userModel models.UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&userModel).Error; err != nil {
		return nil, "", err
	}
	return &domain.User{
		ID:        userModel.ID,
		Username:  userModel.Username,
		FullName:  userModel.FullName,
		Email:     userModel.Email,
		CreatedAt: userModel.CreatedAt,
	}, userModel.HashedPassword, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	var userModel models.UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		return nil, "", err
	}
	return &domain.User{
		ID:        userModel.ID,
		Username:  userModel.Username,
		FullName:  userModel.FullName,
		Email:     userModel.Email,
		CreatedAt: userModel.CreatedAt,
	}, userModel.HashedPassword, nil
}
