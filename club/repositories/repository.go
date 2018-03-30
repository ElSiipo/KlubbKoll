package repository

import models "github.com/klubbkoll/club"

type ClubRepository interface {
	GetByID(id int64) (*models.Club, error)
	GetByTitle(title string) (*models.Club, error)
	Update(club *models.Club) (*models.Club, error)
	Store(a *models.Club) (int64, error)
	Delete(id int64) (bool, error)
}
