package repository

import (
	"fiberStore/models"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) models.ProductRepository {
	return &productRepository{db}
}

func (pr *productRepository) InsertOne(req *models.Product) (*models.Product, error) {
	err := pr.db.Create(req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (pr *productRepository) FindOne(id int) (*models.Product, error) {
	var product *models.Product

	err := pr.db.Model(&product).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *productRepository) FindOneBySlug(slug string) (*models.Product, error) {
	var product *models.Product

	err := pr.db.Model(&product).Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *productRepository) FindByCategory(page, limit int, search string) (*[]models.Product, int, error) {
	var (
		products []models.Product
		count    int64
	)

	err := pr.db.Model(models.Product{}).Where("category LIKE ?", "%"+search+"%").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = pr.db.Where("category LIKE ?", "%"+search+"%").Order("created_at DESC").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, int(count), err
	}

	return &products, int(count), nil
}

func (pr *productRepository) FindAll(page, limit int) (*[]models.Product, int, error) {
	var (
		product  *models.Product
		products []models.Product
		count    int64
	)

	err := pr.db.Model(&product).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = pr.db.Model(&product).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return &products, int(count), nil
}

func (pr *productRepository) FindQueryAll(page, limit int, search string) (*[]models.Product, int, error) {
	var (
		products []models.Product
		count    int64
	)

	err := pr.db.Model(models.Product{}).Where("name LIKE ?", "%"+search+"%").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = pr.db.Where("name LIKE ?", "%"+search+"%").Order("created_at DESC").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, int(count), err
	}

	return &products, int(count), nil
}

func (pr *productRepository) UpdateOne(req *models.Product) (*models.Product, error) {
	err := pr.db.Save(req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (pr *productRepository) DeleteOne(product *models.Product) error {
	err := pr.db.Delete(product).Error
	if err != nil {
		return err
	}

	return nil
}
