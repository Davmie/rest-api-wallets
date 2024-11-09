package postgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Davmie/javaCode/internal/testBuilders"
	walletRep "github.com/Davmie/javaCode/internal/wallet/repository"
	"github.com/Davmie/javaCode/models"
	"github.com/Davmie/javaCode/pkg/logger"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type WalletRepoTestSuite struct {
	suite.Suite
	db            *sql.DB
	gormDB        *gorm.DB
	mock          sqlmock.Sqlmock
	repo          walletRep.WalletRepositoryI
	walletBuilder *testBuilders.WalletBuilder
}

func TestWalletRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(WalletRepoTestSuite))
}

func (s *WalletRepoTestSuite) BeforeEach(t provider.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("error while creating sql mock")
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatal("error gorm open")
	}

	var logger logger.Logger

	s.db = db
	s.gormDB = gormDB
	s.mock = mock

	s.repo = New(logger, gormDB)
	s.walletBuilder = testBuilders.NewWalletBuilder()
}

func (s *WalletRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *WalletRepoTestSuite) TestCreateWallet(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "wallet" ("uid","amount","id") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(wallet.UID, wallet.Amount, wallet.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectCommit()

	err := s.repo.Create(&wallet)
	t.Assert().NoError(err)
	t.Assert().Equal(1, wallet.ID)
}

func (s *WalletRepoTestSuite) TestGetWallet(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("").
		WithAmount(20).
		Build()

	rows := sqlmock.NewRows([]string{"id", "wallet_uid", "amount"}).
		AddRow(
			wallet.ID,
			wallet.UID,
			wallet.Amount,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "wallet" WHERE id = $1 LIMIT $2`)).
		WithArgs(wallet.ID, 1).
		WillReturnRows(rows)

	resWallet, err := s.repo.Get(wallet.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(wallet, *resWallet)
}

func (s *WalletRepoTestSuite) TestUpdateWallet(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	rows := sqlmock.NewRows([]string{"id", "wallet_uid", "amount"}).
		AddRow(
			wallet.ID,
			wallet.UID,
			wallet.Amount,
		)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE "wallet" SET "uid"=$1,"amount"=$2 WHERE "id" = $3 RETURNING *`)).
		WithArgs(wallet.UID, wallet.Amount, wallet.ID).WillReturnRows(rows)

	s.mock.ExpectCommit()

	err := s.repo.Update(&wallet)
	t.Assert().NoError(err)
}

func (s *WalletRepoTestSuite) TestDeleteWallet(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("uid").
		WithAmount(20).
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "wallet" WHERE "wallet"."id" = $1`)).
		WithArgs(wallet.ID).WillReturnResult(sqlmock.NewResult(int64(wallet.ID), 1))

	s.mock.ExpectCommit()

	err := s.repo.Delete(wallet.ID)
	t.Assert().NoError(err)
}

func (s *WalletRepoTestSuite) TestGetAll(t provider.T) {
	wallets := make([]models.Wallet, 10)
	for _, wallet := range wallets {
		err := faker.FakeData(&wallet)
		t.Assert().NoError(err)
	}

	walletsPtr := make([]*models.Wallet, len(wallets))
	for i, wallet := range wallets {
		walletsPtr[i] = &wallet
	}

	rowsWallets := sqlmock.NewRows([]string{"id", "wallet_uid", "amount"})

	for i := range wallets {
		rowsWallets.AddRow(wallets[i].ID, wallets[i].UID, wallets[i].Amount)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "wallet"`)).
		WillReturnRows(rowsWallets)

	resWallets, err := s.repo.GetAll()
	t.Assert().NoError(err)
	t.Assert().Equal(walletsPtr, resWallets)
}

func (s *WalletRepoTestSuite) TestGetByUID(t provider.T) {
	wallet := s.walletBuilder.
		WithID(1).
		WithUID("").
		WithAmount(20).
		Build()

	rows := sqlmock.NewRows([]string{"id", "wallet_uid", "amount"}).
		AddRow(
			wallet.ID,
			wallet.UID,
			wallet.Amount,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "wallet" WHERE uid = $1 LIMIT $2`)).
		WithArgs(wallet.UID, 1).
		WillReturnRows(rows)

	resWallet, err := s.repo.GetByUID(wallet.UID)
	t.Assert().NoError(err)
	t.Assert().Equal(wallet, *resWallet)
}
