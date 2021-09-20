package store_test

import (
	"context"
	"github.com/dsthakur2711/wallet/model"
	"github.com/dsthakur2711/wallet/pkg/local_errors"
	"github.com/dsthakur2711/wallet/store"
	"github.com/dsthakur2711/wallet/store/mocks"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	logs "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)


var testDB *gorm.DB
func DbConn() (db *gorm.DB) {
	/*dbDriver := "mysql"
	dbUser := "root"
	dbPass := "password"
	dbName := "walletDB"*/
	//db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/walletDB")
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/walletDB?parseTime=true")
	logs.Print("db connection opened")
	if err != nil {
		panic(err.Error())
	}
	return db
}

////var mocktransRepo= 	mocks.TransRepo{}
//var transRepo = NewTransRepo(testDB)
//var walletRepoImpl= NewWalletRepo(testDB, transRepo)


func InitWalletRepo(t *testing.T) store.WalletRepo{
	testDB 		 = DbConn()
	transRepo	:= store.NewTransRepo(testDB)
	walletRepo	:= store.NewWalletRepo(testDB, transRepo)

	require.NotEmpty(t, transRepo)
	require.NotEmpty(t, walletRepo)
	return walletRepo
}

func TestWalletRepository_GetWalletByAddressErrWalletNotFound(t *testing.T) {

	walletRepoImpl := InitWalletRepo(t)

	i, err := walletRepoImpl.GetWalletByAddress(context.Background(), "walletAddress")

	//assertion
	assert.Error(t, local_errors.ErrWalletNotFound, err)
	assert.Equal(t, model.Wallet{}, i)
}

func TestWalletRepository_GetWalletByUsernameErrWalletNotFound(t *testing.T) {

	walletRepoImpl := InitWalletRepo(t)

	i, err := walletRepoImpl.GetWalletByUsername(context.Background(), "thisUsernameDoesNotExist")

	//assertion
	assert.Error(t, local_errors.ErrWalletNotFound, err)
	assert.Equal(t, model.Wallet{}, i)
}

func TestWalletRepository_UpdateWalletStatusFailOnErrSomethingWrong(t *testing.T) {
	walletRepoImpl := InitWalletRepo(t)

	i, err := walletRepoImpl.UpdateWalletStatus(context.Background(), store.UpdateWalletStatusParams{})
	var i1 model.Wallet
	i1.UpdatedAt = i.UpdatedAt

	//assertion
	assert.Error(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, i1, i)
}

func TestWalletRepository_CreateWalletFailOnErrUserNotFound(t *testing.T) {
	walletRepoImpl := InitWalletRepo(t)

	i, err := walletRepoImpl.CreateWallet(context.Background(), store.CreateWalletParams{
		Username: "thisUsernameDoesntExist",
		Currency: "INR",
	})
	var i1 model.Wallet
	i1.UpdatedAt = i.UpdatedAt

	//assertion
	assert.Error(t, local_errors.ErrUserNotFound, err)
	assert.Equal(t, i1, i)
}


var mocktransRepo= 	mocks.TransRepo{}
//func TestWalletRepository_SendMoneyFailOnCreateTransfer(t *testing.T) {
//	walletRepoImpl := InitWalletRepo(t)
//
//	//line 160 we need a mock
//	mocktransRepo.On("CreateTransfer", mock.Anything, mock.Anything).
//		Return(model.Trans{}, local_errors.ErrSomethingWrong).Once()
//
//	walletTransferResult, err := walletRepoImpl.SendMoney(context.Background(), store.SendMoneyParams{})
//
//	//assertion
//	assert.Error(t, local_errors.ErrUserNotFound, err)
//	assert.Equal(t, store.WalletTransferResult{}, walletTransferResult)
//}
//
//func TestWalletRepository_SendMoneyFailOnAddMoney(t *testing.T) {
//		walletRepoImpl := InitWalletRepo(t)
//
//		//line 160 we need a mock
//		mocktransRepo.On("CreateTransfer", mock.Anything, mock.Anything).
//			Return(model.Trans{}, nil).Once()
//
//		//line 169 we need a mock
//
//		walletTransferResult, err := walletRepoImpl.SendMoney(context.Background(), store.SendMoneyParams{})
//
//		//assertion
//		assert.Error(t, local_errors.ErrUserNotFound, err)
//		assert.Equal(t, store.WalletTransferResult{}, walletTransferResult)
//}

func TestWalletRepository_AddWalletBalanceErrSomethingWrong(t *testing.T) {
	walletRepoImpl := InitWalletRepo(t)

	i, err := walletRepoImpl.AddWalletBalance(context.Background(), store.AddWalletBalanceParams{})
	var i1 model.Wallet
	i1.UpdatedAt = i.UpdatedAt

	//assertion
	assert.Error(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, i1, i)
}

var mockWalletRepo = mocks.WalletRepo{}
func TestWalletRepository_AddMoney(t *testing.T) {
	walletRepoImpl := InitWalletRepo(t)

	//line 200 we need a mock
	mockWalletRepo.On("AddWalletBalance", mock.Anything, mock.Anything).
		Return(model.Wallet{}, nil).Once()

	aw1, aw2, err := walletRepoImpl.AddMoney(context.Background(), "address1",  5, "address2", -5)

	var ew1, ew2 model.Wallet
	ew1.Balance= 	5
	ew1.UpdatedAt= aw1.UpdatedAt
	ew2.Balance=   -5
	ew2.UpdatedAt= aw2.UpdatedAt
	//assertion
	assert.Nil(t, err)
	assert.Equal(t, ew1, aw1)
	assert.Equal(t, ew2, aw2)
}