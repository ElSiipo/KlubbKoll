package repository

import models "github.com/ElSiipo/klubbkoll/club"

type ClubRepository interface {
	GetAll() ([]*models.Club, error)
	GetByID(id int64) (*models.Club, error)
	Update(club *models.Club) (*models.Club, error)
	Store(a *models.Club) (int64, error)
	Delete(id int64) (bool, error)
}
