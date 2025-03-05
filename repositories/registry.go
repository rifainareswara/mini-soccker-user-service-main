package repositories

import (
	"gorm.io/gorm"
	repositories "user-service/repositories/user"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() repositories.IUserRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetUser() repositories.IUserRepository {
	return repositories.NewUserRepository(r.db)
}
