package repository

import "github.com/Davmie/javaCode/models"

type WalletRepositoryI interface {
	Create(w *models.Wallet) error
	Get(id int) (*models.Wallet, error)
	Update(w *models.Wallet) error
	Delete(id int) error
	GetAll() ([]*models.Wallet, error)
	GetByUID(uid string) (*models.Wallet, error)
}
