package service

import (
	"context"
	"fmt"
	"github.com/dsthakur2711/wallet/dto"
	"github.com/dsthakur2711/wallet/model"
	"github.com/dsthakur2711/wallet/pkg/local_errors"
	"github.com/dsthakur2711/wallet/store"
	"github.com/dsthakur2711/wallet/store/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var mockWalletRepo= mocks.WalletRepo{}
var walletServiceImpl= NewWalletService(&mockWalletRepo)

func TestAddWalletFail(t *testing.T){
	//line 43  we need a mock
	mockWalletRepo.On("CreateWallet", mock.Anything, mock.Anything).
		Return(model.Wallet{}, fmt.Errorf("error")).Once()

	walletDto, err := walletServiceImpl.AddWallet(context.Background(),dto.CreateWalletDto{})

	//assertion
	assert.Error(t, err)
	assert.Equal(t, dto.WalletDto{},walletDto)
}

func TestAddWalletSuccess(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "wa.String()",
		Status: model.WalletStatusACTIVE,
		Balance: 0,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	mockWalletRepo.On("CreateWallet", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	walletDto, err := walletServiceImpl.AddWallet(context.Background(), dto.CreateWalletDto{
		Username: w.Username,
		Currency: w.Currency,
	})

	//assertion
	assert.Nil(t, err)
	assert.Equal(t, dto.NewWalletDto(w),walletDto)
}

func TestWalletService_GetWalletFailErrWalletNotFound(t *testing.T) {
	//line 77 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrWalletNotFound).Once()

	walletDto, err := walletServiceImpl.GetWallet(context.Background(), dto.GetWalletDto{})

	//assertion
	assert.Equal(t, local_errors.ErrWalletNotFound, err)
	assert.Equal(t, dto.WalletDto{}, walletDto)
}

func TestWalletService_GetWalletFailErrSomethingWrong(t *testing.T) {
	//line 77 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrSomethingWrong).Once()

	walletDto, err := walletServiceImpl.GetWallet(context.Background(), dto.GetWalletDto{})

	//assertion
	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.WalletDto{}, walletDto)
}

func TestWalletService_GetWalletSuccess(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "wa.String()",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "Currency",
		//CreatedAt: time.Now(),
	}
	//line 77 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	walletDto, err := walletServiceImpl.GetWallet(context.Background(), dto.GetWalletDto{})

	//assertion
	assert.Nil(t, err)
	assert.Equal(t, dto.NewWalletDto(w), walletDto)
}

func TestWalletService_CreditInvalidAmount(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "wa.String()",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}

	updatedWalletBalanceDto, err := walletServiceImpl.Credit(context.Background(), dto.CreditDto{
		WalletAddress: w.WalletAddress,
		Amount: -3,
	})

	//assertion
	assert.Error(t, err)
	assert.Equal(t, dto.UpdatedWalletBalanceDto{}, updatedWalletBalanceDto)
}

func TestWalletService_CreditFailWalletAddressDoesNotExist(t *testing.T) {
	//line 154 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrWalletNotFound).Once()

	updatedWalletBalanceDto, err := walletServiceImpl.Credit(context.Background(), dto.CreditDto{
		WalletAddress: "ThisIsWalletAddress",
		Amount: 5,
	})

	//assertion
	assert.Equal(t, local_errors.ErrWalletNotFound, err)
	assert.Equal(t, dto.UpdatedWalletBalanceDto{}, updatedWalletBalanceDto)
}

func TestWalletService_CreditGetWalletErrSomethingWrong(t *testing.T) {
	//line 154 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrSomethingWrong).Once()

	updatedWalletBalanceDto, err := walletServiceImpl.Credit(context.Background(), dto.CreditDto{
		WalletAddress: "ThisIsWalletAddress",
		Amount: 5,
	})

	//assertion
	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.UpdatedWalletBalanceDto{}, updatedWalletBalanceDto)
}


func TestWalletService_CreditFailOnAddWalletBalance(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "wa.String()",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	//line 154 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()
	//line 160 we need a mock
	mockWalletRepo.On("AddWalletBalance", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrSomethingWrong).Once()

	updatedWalletBalanceDto, err := walletServiceImpl.Credit(context.Background(), dto.CreditDto{
		WalletAddress: w.WalletAddress,
		Amount: 3,
	})

	//assertion
	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.UpdatedWalletBalanceDto{}, updatedWalletBalanceDto)
}

func TestWalletService_CreditSuccess(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "wa.String()",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	w1 := model.Wallet{
		Username: "deepak",
		WalletAddress: "wa.String()",
		Status: model.WalletStatusACTIVE,
		Balance: 26,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	//line 154 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()
	//line 160 we need a mock
	mockWalletRepo.On("AddWalletBalance", mock.Anything, mock.Anything).
		Return(w1, nil).Once()

	updatedWalletBalanceDto, err := walletServiceImpl.Credit(context.Background(), dto.CreditDto{
		WalletAddress: w.WalletAddress,
		Amount: 3,
	})

	//assertion
	assert.Nil(t, err)
	assert.Equal(t, dto.NewUpdatedWalletBalanceDto(w1), updatedWalletBalanceDto)
}

func TestWalletService_PayFromWalletErr(t *testing.T) {
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrWalletNotFound).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Equal(t, fmt.Errorf("from_wallet_address does not exists"), err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}

func TestWalletService_PayFromWalletErrSomethingWrong(t *testing.T) {
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrSomethingWrong).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}

func TestWalletService_PayFromWalletStatusInactive(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusINACTIVE,
		Balance: 23,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Equal(t, fmt.Errorf("inactive from_wallet"), err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}


func TestWalletService_PayToWalletErr(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()


	//line 116 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrWalletNotFound).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Equal(t, fmt.Errorf("to_wallet_address does not exists"), err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}

func TestWalletService_PayToWalletErrSomethingWrong(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	//line 116 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrSomethingWrong).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}

func TestWalletService_PayToWalletStatusInactive(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	w1 := model.Wallet{
		Username: "deepak",
		WalletAddress: "toWallet",
		Status: model.WalletStatusINACTIVE,
		Balance: 23,
		Currency: "Currency",
		CreatedAt: time.Now(),
	}
	//line 116 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w1, nil).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Equal(t, fmt.Errorf("inactive to_wallet"), err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}

func TestWalletService_PayFailCurrencyMisMatch(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	w1 := model.Wallet{
		Username: "deepak",
		WalletAddress: "toWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "USD",
		CreatedAt: time.Now(),
	}
	//line 116 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w1, nil).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Equal(t, fmt.Errorf("Currency of both wallets should be same"), err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}

func TestWalletService_PayFailAmountNotPositive(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	w1 := model.Wallet{
		Username: "deepak",
		WalletAddress: "toWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 116 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w1, nil).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			-3,
	})

	assert.Equal(t, fmt.Errorf("amount to pay should be positive"), err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}

func TestWalletService_PayFailBalanceNotSufficientInFromWallet(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 2,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	w1 := model.Wallet{
		Username: "deepak",
		WalletAddress: "toWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 116 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w1, nil).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Equal(t, fmt.Errorf("insufficient wallet balance"), err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}


func TestWalletService_PayFailSendMoney(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 25,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	w1 := model.Wallet{
		Username: "deepak",
		WalletAddress: "toWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 116 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w1, nil).Once()

	//line 142 we need a mock
	mockWalletRepo.On("SendMoney", mock.Anything, mock.Anything).
		Return(store.WalletTransferResult{}, fmt.Errorf("error")).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Error(t, err)
	assert.Equal(t, dto.TransResultDto{}, transResultDto)
}

func TestWalletService_PaySuccess(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "fromWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 25,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 101 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	w1 := model.Wallet{
		Username: "deepak",
		WalletAddress: "toWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 116 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w1, nil).Once()

	trans := model.Trans{
		FromWalletAdd: "fromWalletAddress",
		ToWalletAdd: 	"toWalletAddress",
		Amount: 		3,
		CreatedAt: 		time.Now(),
	}
	//line 142 we need a mock
	mockWalletRepo.On("SendMoney", mock.Anything, mock.Anything).
		Return(store.WalletTransferResult{Trans: trans}, nil).Once()

	transResultDto, err := walletServiceImpl.Pay(context.Background(), dto.TransferMoneyDto{
		FromWalletAddress: "fromWalletAddress",
		ToWalletAddress:   "toWalletAddress",
		Amount: 			3,
	})

	assert.Nil(t, err)
	assert.Equal(t, dto.NewTransResultDto(trans), transResultDto)
}

func TestWalletService_UpdateWalletStatusFailOnErrWalletNotFound(t *testing.T) {
	//line 188 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrWalletNotFound).Once()

	walletStatusDto, err := walletServiceImpl.UpdateWalletStatus(context.Background(), dto.UpdateWalletStatusDto{})

	assert.Equal(t, fmt.Errorf("wallet_address does not exists"), err)
	assert.Equal(t, dto.UpdateWalletStatusDto{}, walletStatusDto)
}

func TestWalletService_UpdateWalletStatusFailOnErrSomethingWrong(t *testing.T) {
	//line 188 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrSomethingWrong).Once()

	walletStatusDto, err := walletServiceImpl.UpdateWalletStatus(context.Background(), dto.UpdateWalletStatusDto{})

	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.UpdateWalletStatusDto{}, walletStatusDto)
}

func TestWalletService_UpdateWalletStatusAlreadySameStatus(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "aWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 188 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	walletStatusDto, err := walletServiceImpl.UpdateWalletStatus(context.Background(), dto.UpdateWalletStatusDto{
		WalletAddress: 	"aWallet",
		Status: 		model.WalletStatusACTIVE,
	})

	assert.Nil(t, err)
	assert.Equal(t, dto.NewUpdateWalletStatusDto(w), walletStatusDto)
}

func TestWalletService_UpdateWalletStatusInvalidStatusChangeRequest(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "aWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 188 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	walletStatusDto, err := walletServiceImpl.UpdateWalletStatus(context.Background(), dto.UpdateWalletStatusDto{
		WalletAddress: 	"aWallet",
		Status: 		"ThisIsInvalidStatusChangeRequest",
	})

	assert.Equal(t, fmt.Errorf("Invalid status change request"), err)
	assert.Equal(t, dto.UpdateWalletStatusDto{}, walletStatusDto)
}


func TestWalletService_UpdateWalletStatusErr(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "aWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 188 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	//line 210 we need a mock
	mockWalletRepo.On("UpdateWalletStatus", mock.Anything, mock.Anything).
		Return(model.Wallet{}, local_errors.ErrSomethingWrong).Once()

	walletStatusDto, err := walletServiceImpl.UpdateWalletStatus(context.Background(), dto.UpdateWalletStatusDto{
		WalletAddress: 	"aWallet",
		Status: 		"INACTIVE",
	})

	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.UpdateWalletStatusDto{}, walletStatusDto)
}


func TestWalletService_UpdateWalletStatusSuccess(t *testing.T) {
	w := model.Wallet{
		Username: "deepak",
		WalletAddress: "aWallet",
		Status: model.WalletStatusACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	w1 := model.Wallet{
		Username: "deepak",
		WalletAddress: "aWallet",
		Status: model.WalletStatusINACTIVE,
		Balance: 23,
		Currency: "INR",
		CreatedAt: time.Now(),
	}
	//line 188 we need a mock
	mockWalletRepo.On("GetWalletByAddress", mock.Anything, mock.Anything).
		Return(w, nil).Once()

	//line 210 we need a mock
	mockWalletRepo.On("UpdateWalletStatus", mock.Anything, mock.Anything).
		Return(w1, nil)

	walletStatusDto, err := walletServiceImpl.UpdateWalletStatus(context.Background(), dto.UpdateWalletStatusDto{
		WalletAddress: 	"aWallet",
		Status: 		"INACTIVE",
	})

	assert.Nil(t, err)
	assert.Equal(t, dto.NewUpdateWalletStatusDto(w1), walletStatusDto)
}
