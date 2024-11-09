package postgres

import (
	"github.com/Davmie/javaCode/internal/wallet/repository"
	"github.com/Davmie/javaCode/models"
	"github.com/Davmie/javaCode/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pgWalletRepo struct {
	Logger logger.Logger
	DB     *gorm.DB
}

func New(logger logger.Logger, db *gorm.DB) repository.WalletRepositoryI {
	return &pgWalletRepo{
		Logger: logger,
		DB:     db,
	}
}

func (pr *pgWalletRepo) Create(w *models.Wallet) error {
	tx := pr.DB.Create(w)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgWalletRepo.Create error while inserting in repo")
	}

	return nil
}

func (pr *pgWalletRepo) Get(id int) (*models.Wallet, error) {
	var w models.Wallet
	tx := pr.DB.Where("id = ?", id).Take(&w)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgWalletRepo.Get error")
	}

	return &w, nil
}

func (pr *pgWalletRepo) Update(w *models.Wallet) error {
	tx := pr.DB.Clauses(clause.Returning{}).Omit("id").Select("amount", "uid").Updates(w)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgWalletRepo.Update error while inserting in repo")
	}

	return nil
}

func (pr *pgWalletRepo) Delete(id int) error {
	tx := pr.DB.Delete(&models.Wallet{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgWalletRepo.Delete error")
	}

	return nil
}

func (pr *pgWalletRepo) GetAll() ([]*models.Wallet, error) {
	var wallets []*models.Wallet

	tx := pr.DB.Find(&wallets)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgWalletRepo.GetAll error")
	}

	return wallets, nil
}

func (pr *pgWalletRepo) GetByUID(uid string) (*models.Wallet, error) {
	var w models.Wallet
	tx := pr.DB.Where("uid = ?", uid).Take(&w)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgWalletRepo.Get error")
	}

	return &w, nil
}
