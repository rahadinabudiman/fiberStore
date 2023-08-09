package repository

import (
	"fiberStore/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	InsertOne(req *models.User) (*models.User, error)
	FindOneById(id int) (*models.User, error)
	FindAll(page, limit int, search string) (*[]models.User, int, error)
	UpdateOne(req *models.User) (*models.User, error)
	DeleteOne(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) InsertOne(req *models.User) (*models.User, error) {
	err := ur.db.Create(req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (ur *userRepository) FindOneById(id int) (*models.User, error) {
	var user *models.User

	err := ur.db.Model(&user).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) FindAll(page, limit int, search string) (*[]models.User, int, error) {
	var users *[]models.User
	var count int64

	err := ur.db.Unscoped().Find(&users).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = ur.db.Unscoped().Where("role = 'user' AND name LIKE ? OR role = 'user' AND username LIKE ?", "%"+search+"%", "%"+search+"%").Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return users, int(count), nil
	}

	return nil, 0, nil
}

func (ur *userRepository) UpdateOne(req *models.User) (*models.User, error) {
	err := ur.db.Save(req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (ur *userRepository) DeleteOne(id int) error {
	err := ur.db.Where("id = ?", id).Delete(&models.User{}).Error
	if err != nil {
		return err
	}

	return nil
}
