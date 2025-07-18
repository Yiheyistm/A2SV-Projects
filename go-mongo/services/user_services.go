package services

import "github.com/yiheyistm/go-mongo/models"

type UserServices interface {
	Create(*models.User) error
	GetAll() ([]models.User, error)
	GetSingle(*string) (models.User, error)
	Update(*models.User) error
	Delete(*string) error
}
