package repository

import models "github.com/ElSiipo/Klubbkoll/club"

type ClubRepository interface {
	GetAll() ([]*models.Club, error)
	GetByID(clubID string) (*models.Club, error)
	Update(club *models.Club) (*models.Club, error)
	Store(a *models.Club) (string, error)
	Delete(clubID string) (bool, error)
}
