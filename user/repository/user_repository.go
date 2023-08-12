package repository

import (
	"fiberStore/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) models.UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) InsertOne(req *models.User) (*models.User, error) {
	err := ur.db.Create(req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (ur *userRepository) FindOneByUsername(username string) (*models.User, error) {
	var user *models.User

	err := ur.db.Model(&user).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) FindOneAdmin(id uint) (*models.User, error) {
	var user *models.User

	err := ur.db.Model(&user).Where("id = ? AND role = 'Admin'", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
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
	var users []models.User
	var count int64

	err := ur.db.Unscoped().Model(&models.User{}).Where("role = 'Customer' AND (name LIKE ? OR username LIKE ?)", "%"+search+"%", "%"+search+"%").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = ur.db.Unscoped().Where("role = 'Customer' AND (name LIKE ? OR username LIKE ?)", "%"+search+"%", "%"+search+"%").Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, int(count), err
	}

	return &users, int(count), nil
}

func (ur *userRepository) UpdateOne(req *models.User) (*models.User, error) {
	err := ur.db.Save(req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (ur *userRepository) DeleteOne(user *models.User) error {
	err := ur.db.Delete(&user).Error
	if err != nil {
		return err
	}

	return nil
}
