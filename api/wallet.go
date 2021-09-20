package api

import (
	"encoding/json"
	"github.com/dsthakur2711/wallet/dto"
	types "github.com/dsthakur2711/wallet/pkg/local_errors"
	"github.com/dsthakur2711/wallet/service"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"net/http"
)

type WalletResource interface {
	AddWallet(w http.ResponseWriter, r *http.Request)
	Pay(w http.ResponseWriter, r *http.Request)
	Credit(w http.ResponseWriter, r *http.Request)
	UpdateWalletStatus(w http.ResponseWriter, r *http.Request)
	GetWallet(w http.ResponseWriter, r *http.Request)
	//RegisterRoutes(r chi.Router)
}

type walletResource struct {
	walletSvc service.WalletSvc
}

func NewWalletResource(walletSvc service.WalletSvc) WalletResource {
	return &walletResource{
		walletSvc: walletSvc,
	}
}

//func (wr *walletResource) RegisterRoutes(r chi.Router) {
//	r.Get("/wallets/{walletID}", wr.Get)
//	r.Post("/wallets/pay", wr.Pay)
//}

func (wr *walletResource) AddWallet(w http.ResponseWriter, r *http.Request) {
	logrus.Println("log AddWallet in api/wallet/AddWallet ")
	var req dto.CreateWalletDto
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}

	wallet, err := wr.walletSvc.AddWallet(ctx, req)
	if err != nil {
		_ = render.Render(w, r, types.ErrResponse(err))
		return
	}

	render.JSON(w, r, wallet)
}


func (wr *walletResource) Pay(w http.ResponseWriter, r *http.Request) {

	logrus.Println("log Pay in api/wallet/Pay ")

	var req dto.TransferMoneyDto
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}

	res, err := wr.walletSvc.Pay(ctx, req)
	if err != nil {
		_ = render.Render(w, r, types.ErrResponse(err))
		return
	}

	render.JSON(w, r, res)
}

func (wr *walletResource) Credit(w http.ResponseWriter, r *http.Request) {

	logrus.Println("log Credit in api/wallet/Credit ")

	var req dto.CreditDto
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}

	res, err := wr.walletSvc.Credit(ctx, req)
	if err != nil {
		_ = render.Render(w, r, types.ErrResponse(err))
		return
	}

	render.JSON(w, r, res)
}

func (wr* walletResource) UpdateWalletStatus(w http.ResponseWriter, r *http.Request) {
	logrus.Println("log UpdateWalletStatus in api/wallet/UpdateWalletStatus ")

	var req dto.UpdateWalletStatusDto
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}

	res, err := wr.walletSvc.UpdateWalletStatus(ctx, req)
	if err != nil {
		_ = render.Render(w, r, types.ErrResponse(err))
		return
	}

	render.JSON(w, r, res)
}

func (wr* walletResource) GetWallet(w http.ResponseWriter, r *http.Request) {
	logrus.Println("log GetWallet in api/wallet/GetWallet ")

	var req dto.GetWalletDto
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(req); err != nil {
		_ = render.Render(w, r, types.ErrBadRequest(err))
		return
	}

	res, err := wr.walletSvc.GetWallet(ctx, req)
	if err != nil {
		_ = render.Render(w, r, types.ErrResponse(err))
		return
	}

	render.JSON(w, r, res)
}