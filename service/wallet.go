package service

import (
	"context"
	"fmt"
	"github.com/dsthakur2711/wallet/dto"
	"github.com/dsthakur2711/wallet/model"
	"github.com/dsthakur2711/wallet/pkg/local_errors"
	"github.com/dsthakur2711/wallet/store"
	"github.com/sirupsen/logrus"
)

type WalletSvc interface {
	Pay(ctx context.Context, transferMoneyDto dto.TransferMoneyDto) (dto.TransResultDto, error)
	Credit(ctx context.Context, creditDto dto.CreditDto) (dto.UpdatedWalletBalanceDto, error)
	UpdateWalletStatus(ctx context.Context, updateWalletStatusDto dto.UpdateWalletStatusDto) (dto.UpdateWalletStatusDto, error)
	AddWallet(ctx context.Context,createWalletDto dto.CreateWalletDto) (dto.WalletDto, error)
	//GetWalletByUsername(ctx context.Context, username string) (dto.WalletDto, error)
	GetWallet(ctx context.Context, getWalletDto dto.GetWalletDto) (dto.WalletDto, error)
}

type walletService struct {
	walletRepo store.WalletRepo
}

func NewWalletService(walletRepo store.WalletRepo) WalletSvc {
	return &walletService{
		walletRepo: walletRepo,
	}
}


func (w *walletService) AddWallet(ctx context.Context,createWalletDto dto.CreateWalletDto) (dto.WalletDto,error){
	logrus.Println("log AddWallet in service/wallet/AddWallet ")

	var walletDto dto.WalletDto

	arg := store.CreateWalletParams{
		Username:	createWalletDto.Username,
		Currency: 	createWalletDto.Currency,
	}

	wallet, err := w.walletRepo.CreateWallet(ctx,arg)
	if err != nil{
		return walletDto, err
	}

	walletDto = dto.NewWalletDto(wallet)
	return walletDto,nil
}

//func (w *walletService) GetWalletByUsername(ctx context.Context, username string) (dto.WalletDto, error) {
//
//	logrus.Println("log GetWalletById in service/wallet/GetWalletById ")
//
//	var walletDto dto.WalletDto
//
//	wallet, err := w.walletRepo.GetWalletByUsername(ctx, username)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return walletDto, local_errors.ErrWalletNotFound
//		}
//
//		return walletDto, err
//	}
//
//	walletDto = dto.NewWalletDto(wallet)
//	return walletDto, nil
//}

func (w *walletService) GetWallet(ctx context.Context, getWalletDto dto.GetWalletDto) (dto.WalletDto, error) {

	logrus.Println("log GetWalletByAddress in service/wallet/GetWalletByAddress ")

	var walletDto dto.WalletDto

	wallet, err := w.walletRepo.GetWalletByAddress(ctx, getWalletDto.WalletAddress)
	if err != nil {
		if err == local_errors.ErrWalletNotFound {
			return walletDto, local_errors.ErrWalletNotFound
		}

		return walletDto, local_errors.ErrSomethingWrong
	}

	walletDto = dto.NewWalletDto(wallet)
	return walletDto, nil
}

func (w *walletService) Pay(ctx context.Context, transferMoneyDto dto.TransferMoneyDto) (dto.TransResultDto, error) {
	logrus.Println("log Pay in service/wallet/Pay ")

	var txnResDto dto.TransResultDto

	arg := store.SendMoneyParams{
		FromWalletAddress: transferMoneyDto.FromWalletAddress,
		ToWalletAddress:   transferMoneyDto.ToWalletAddress,
		Amount:            transferMoneyDto.Amount,
	}

	fromWallet, err := w.walletRepo.GetWalletByAddress(ctx, arg.FromWalletAddress)

	if err != nil {
		if err == local_errors.ErrWalletNotFound {
			return txnResDto, fmt.Errorf("from_wallet_address does not exists")
		}
		return txnResDto, local_errors.ErrSomethingWrong
	}

	if fromWallet.Status != model.WalletStatusACTIVE {
		logrus.Println("log  fromWallet.Status is not ACTIVE !! ")
		return txnResDto, fmt.Errorf("inactive from_wallet")
	}


	toWallet, err := w.walletRepo.GetWalletByAddress(ctx,arg.ToWalletAddress)

	if err != nil {
		if err == local_errors.ErrWalletNotFound {
			return txnResDto, fmt.Errorf("to_wallet_address does not exists")
		}
		return txnResDto, local_errors.ErrSomethingWrong
	}

	if toWallet.Status != model.WalletStatusACTIVE{
		logrus.Println("log  toWallet.Status is not ACTIVE !! ")
		return txnResDto, fmt.Errorf("inactive to_wallet")
	}

	if fromWallet.Currency != toWallet.Currency{
		return txnResDto, fmt.Errorf("Currency of both wallets should be same")
	}

	if arg.Amount <= 0 {
		return txnResDto, fmt.Errorf("amount to pay should be positive")
	}

	if !fromWallet.IsBalanceSufficient(arg.Amount) {
		return txnResDto, fmt.Errorf("insufficient wallet balance")
	}

	res, err := w.walletRepo.SendMoney(ctx, arg)
	if err != nil {
		return txnResDto, err
	}

	txnResDto = dto.NewTransResultDto(res.Trans)
	return txnResDto, nil
}

func (w *walletService) Credit(ctx context.Context, creditDto dto.CreditDto) (dto.UpdatedWalletBalanceDto,error){
	logrus.Println("log Credit in service/wallet/Credit ")

	var updatedWalletBalanceDto dto.UpdatedWalletBalanceDto

	if creditDto.Amount <= 0 {
		return updatedWalletBalanceDto, fmt.Errorf("amount to credit should be positive")
	}

	wallet, err := w.walletRepo.GetWalletByAddress(ctx, creditDto.WalletAddress)
	if err != nil {
		if err == local_errors.ErrWalletNotFound{
			return updatedWalletBalanceDto, err
		}
		return updatedWalletBalanceDto, local_errors.ErrSomethingWrong
	}

	arg := store.AddWalletBalanceParams{
		WalletAddress: 	 creditDto.WalletAddress,
		Amount:			 creditDto.Amount + wallet.Balance,
	}

	wallet1, err := w.walletRepo.AddWalletBalance(ctx,arg)
	if err != nil{
		return updatedWalletBalanceDto, err
	}

	updatedWalletBalanceDto = dto.NewUpdatedWalletBalanceDto(wallet1)
	return updatedWalletBalanceDto,nil
}

func (w *walletService) UpdateWalletStatus(ctx context.Context, updateWalletStatusDto dto.UpdateWalletStatusDto) (dto.UpdateWalletStatusDto, error){
	logrus.Println("log UpdateWalletStatus in service/wallet/CUpdateWalletStatus ")

	var walletStatusDto dto.UpdateWalletStatusDto

	var i model.Wallet
	i, err := w.walletRepo.GetWalletByAddress(ctx, updateWalletStatusDto.WalletAddress)
	if err != nil {
		if err == local_errors.ErrWalletNotFound {
			return walletStatusDto, fmt.Errorf("wallet_address does not exists")
		}
		return walletStatusDto, local_errors.ErrSomethingWrong
	}

	//Already Same Status
	if i.Status == updateWalletStatusDto.Status{
		walletStatusDto = dto.NewUpdateWalletStatusDto(i)
		return walletStatusDto, nil
	}

	if updateWalletStatusDto.Status != model.WalletStatusACTIVE && updateWalletStatusDto.Status != model.WalletStatusINACTIVE{
		return walletStatusDto, fmt.Errorf("Invalid status change request")
	}

	arg := store.UpdateWalletStatusParams{
		WalletAddress: updateWalletStatusDto.WalletAddress,
		Status: 	   updateWalletStatusDto.Status,
	}
	wallet, err := w.walletRepo.UpdateWalletStatus(ctx,arg)
	if err != nil {
		return walletStatusDto, err
	}

	walletStatusDto = dto.NewUpdateWalletStatusDto(wallet)
	return walletStatusDto, nil
}