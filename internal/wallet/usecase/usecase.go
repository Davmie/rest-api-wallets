package usecase

import (
	walletRep "github.com/Davmie/javaCode/internal/wallet/repository"
	"github.com/Davmie/javaCode/models"
	"github.com/pkg/errors"
)

type WalletUseCaseI interface {
	Create(w *models.Wallet) error
	Get(id int) (*models.Wallet, error)
	Update(w *models.Wallet) error
	Delete(id int) error
	GetAll() ([]*models.Wallet, error)
	GetByUID(uid string) (*models.Wallet, error)
	ChangeAmount(uid string, amount int) error
}

type walletUseCase struct {
	walletRepository walletRep.WalletRepositoryI
}

func New(wRep walletRep.WalletRepositoryI) WalletUseCaseI {
	return &walletUseCase{
		walletRepository: wRep,
	}
}

func (wUC *walletUseCase) Create(w *models.Wallet) error {
	err := wUC.walletRepository.Create(w)

	if err != nil {
		return errors.Wrap(err, "walletUseCase.Create error")
	}

	return nil
}

func (wUC *walletUseCase) Get(id int) (*models.Wallet, error) {
	resWallet, err := wUC.walletRepository.Get(id)

	if err != nil {
		return nil, errors.Wrap(err, "walletUseCase.Get error")
	}

	return resWallet, nil
}

func (wUC *walletUseCase) Update(w *models.Wallet) error {
	_, err := wUC.walletRepository.Get(w.ID)

	if err != nil {
		return errors.Wrap(err, "walletUseCase.Update error: Wallet not found")
	}

	err = wUC.walletRepository.Update(w)

	if err != nil {
		return errors.Wrap(err, "walletUseCase.Update error: Can't update in repo")
	}

	return nil
}

func (wUC *walletUseCase) Delete(id int) error {
	_, err := wUC.walletRepository.Get(id)

	if err != nil {
		return errors.Wrap(err, "walletUseCase.Delete error: Wallet not found")
	}

	err = wUC.walletRepository.Delete(id)

	if err != nil {
		return errors.Wrap(err, "walletUseCase.Delete error: Can't delete in repo")
	}

	return nil
}

func (wUC *walletUseCase) GetAll() ([]*models.Wallet, error) {
	wallets, err := wUC.walletRepository.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "walletUseCase.GetAll error")
	}

	return wallets, nil
}

func (wUC *walletUseCase) GetByUID(uid string) (*models.Wallet, error) {
	wallet, err := wUC.walletRepository.GetByUID(uid)
	if err != nil {
		return nil, errors.Wrap(err, "walletUseCase.GetByUID error")
	}

	return wallet, nil
}

func (wUC *walletUseCase) ChangeAmount(uid string, amount int) error {
	wallet, err := wUC.walletRepository.GetByUID(uid)
	if err != nil {
		return errors.Wrap(err, "walletUseCase.ChangeAmount error: Wallet not found")
	}

	wallet.Amount += amount

	err = wUC.walletRepository.Update(wallet)
	if err != nil {
		return errors.Wrap(err, "walletUseCase.ChangeAmount error: Can't update in repo")
	}

	return nil
}
