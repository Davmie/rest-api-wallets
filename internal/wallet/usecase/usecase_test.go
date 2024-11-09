package usecase

import (
	"github.com/Davmie/javaCode/internal/testBuilders"
	walletMocks "github.com/Davmie/javaCode/internal/wallet/repository/mocks"
	"github.com/Davmie/javaCode/models"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"testing"
)

type WalletTestSuite struct {
	suite.Suite
	uc             WalletUseCaseI
	walletRepoMock *walletMocks.WalletRepositoryI
	walletBuilder  *testBuilders.WalletBuilder
}

func TestWalletTestSuite(t *testing.T) {
	suite.RunSuite(t, new(WalletTestSuite))
}

func (s *WalletTestSuite) BeforeEach(t provider.T) {
	s.walletRepoMock = walletMocks.NewWalletRepositoryI(t)
	s.uc = New(s.walletRepoMock)
	s.walletBuilder = testBuilders.NewWalletBuilder()
}

func (s *WalletTestSuite) TestCreateWallet(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	s.walletRepoMock.On("Create", &wallet).Return(nil)
	err := s.uc.Create(&wallet)

	t.Assert().NoError(err)
	t.Assert().Equal(wallet.ID, 1)
}

func (s *WalletTestSuite) TestUpdateWallet(t provider.T) {
	var err error
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	notFoundWallet := s.walletBuilder.WithID(0).Build()

	s.walletRepoMock.On("Get", wallet.ID).Return(&wallet, nil)
	s.walletRepoMock.On("Update", &wallet).Return(nil)
	s.walletRepoMock.On("Get", notFoundWallet.ID).Return(&notFoundWallet, errors.Wrap(err, "Wallet not found"))
	s.walletRepoMock.On("Update", &notFoundWallet).Return(errors.Wrap(err, "Wallet not found"))

	cases := map[string]struct {
		ArgData *models.Wallet
		Error   error
	}{
		"success": {
			ArgData: &wallet,
			Error:   nil,
		},
		"Wallet not found": {
			ArgData: &notFoundWallet,
			Error:   errors.Wrap(err, "Wallet not found"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.Update(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *WalletTestSuite) TestGetWallet(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	s.walletRepoMock.On("Get", wallet.ID).Return(&wallet, nil)
	result, err := s.uc.Get(wallet.ID)

	t.Assert().NoError(err)
	t.Assert().Equal(&wallet, result)
}

func (s *WalletTestSuite) TestDeleteWallet(t provider.T) {
	var err error
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	notFoundWallet := s.walletBuilder.WithID(0).Build()

	s.walletRepoMock.On("Get", wallet.ID).Return(&wallet, nil)
	s.walletRepoMock.On("Delete", wallet.ID).Return(nil)
	s.walletRepoMock.On("Get", notFoundWallet.ID).Return(&notFoundWallet, errors.Wrap(err, "Wallet not found"))
	s.walletRepoMock.On("Delete", notFoundWallet.ID).Return(errors.Wrap(err, "Wallet not found"))

	cases := map[string]struct {
		WalletID int
		Error    error
	}{
		"success": {
			WalletID: wallet.ID,
			Error:    nil,
		},
		"Wallet not found": {
			WalletID: notFoundWallet.ID,
			Error:    errors.Wrap(err, "Wallet not found"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.Delete(test.WalletID)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *WalletTestSuite) TestGetAll(t provider.T) {
	wallets := make([]models.Wallet, 0, 10)
	err := faker.FakeData(&wallets)
	t.Assert().NoError(err)

	walletsPtr := make([]*models.Wallet, len(wallets))
	for i, wallet := range wallets {
		walletsPtr[i] = &wallet
	}

	s.walletRepoMock.On("GetAll").Return(walletsPtr, nil)

	cases := map[string]struct {
		Wallets []models.Wallet
		Error   error
	}{
		"success": {
			Wallets: wallets,
			Error:   nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resWallets, err := s.uc.GetAll()
			t.Assert().ErrorIs(err, test.Error)
			t.Assert().Equal(walletsPtr, resWallets)
		})
	}
}

func (s *WalletTestSuite) TestGetByUID(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	s.walletRepoMock.On("GetByUID", wallet.UID).Return(&wallet, nil)
	result, err := s.uc.GetByUID(wallet.UID)

	t.Assert().NoError(err)
	t.Assert().Equal(&wallet, result)
}

func (s *WalletTestSuite) TestChangeAmount(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	s.walletRepoMock.On("GetByUID", wallet.UID).Return(&wallet, nil)
	wallet.Amount += 1000
	s.walletRepoMock.On("Update", &wallet).Return(nil)

	err := s.uc.ChangeAmount("uid", 1000)

	t.Assert().NoError(err)
}
