package resource

import (
	"errors"
	"gorm.io/gorm"
)

var ErrResourceNameDuplicate = errors.New("resource name already exists")
var ErrResourceNotFound = errors.New("resource not found")

type Repository interface {
	GetResources() ([]Resource, error)
	AddResource(resource *Resource) error
	UpdateResource(resource *Resource) error
	DeleteResource(id string) error
	GetResourceByID(id string) (*Resource, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetResources() ([]Resource, error) {
	var resources []Resource
	result := r.db.Find(&resources)
	return resources, result.Error
}

func (r *GormRepository) AddResource(resource *Resource) error {
	if err := r.db.Create(resource).Error; err != nil {
		// Map DB duplicate key violation to service-domain conflict error.
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrResourceNameDuplicate
		}
		return err
	}

	return nil
}

func (r *GormRepository) UpdateResource(resource *Resource) error {
	result := r.db.Model(&Resource{}).
		Where("id = ?", resource.ID).
		Updates(resource)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrResourceNameDuplicate
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrResourceNotFound
	}
	return nil
}

func (r *GormRepository) DeleteResource(id string) error {
	result := r.db.Delete(&Resource{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrResourceNotFound
	}
	return nil
}

func (r *GormRepository) GetResourceByID(id string) (*Resource, error) {
	var resource Resource
	result := r.db.First(&resource, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, result.Error
	}
	return &resource, nil
}
