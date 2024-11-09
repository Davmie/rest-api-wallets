package delivery

import (
	"encoding/json"
	"fmt"
	walletUseCase "github.com/Davmie/javaCode/internal/wallet/usecase"
	"io"
	"net/http"
	"strconv"

	"github.com/Davmie/javaCode/models"
	"github.com/Davmie/javaCode/pkg/logger"
	"github.com/asaskevich/govalidator"
)

type WalletHandler struct {
	WalletUseCase walletUseCase.WalletUseCaseI
	Logger        logger.Logger
}

func (ah *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {
	wallet := models.Wallet{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &wallet)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(wallet)
	//if err != nil {
	//	ah.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	err = ah.WalletUseCase.Create(&wallet)
	if err != nil {
		ah.Logger.Infow("can`t create wallet",
			"err:", err.Error())
		http.Error(w, "can`t create wallet", http.StatusBadRequest)
		return
	}

	//resp, err := json.Marshal(wallet)
	//
	//if err != nil {
	//	ah.Logger.Errorw("can`t marshal wallet",
	//		"err:", err.Error())
	//	http.Error(w, "can`t make wallet", http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/wallets/%d", wallet.ID))
	w.WriteHeader(http.StatusCreated)

	//_, err = w.Write(resp)
	//if err != nil {
	//	ah.Logger.Errorw("can`t write response",
	//		"err:", err.Error())
	//	http.Error(w, "can`t write response", http.StatusInternalServerError)
	//	return
	//}
}

func (ah *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	walletIdString := r.PathValue("walletId")
	if walletIdString == "" {
		ah.Logger.Errorw("no walletId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	walletId, err := strconv.Atoi(walletIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	wallet, err := ah.WalletUseCase.Get(walletId)
	if err != nil {
		ah.Logger.Infow("can`t get wallet",
			"err:", err.Error())
		http.Error(w, "can`t get wallet", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(wallet)

	if err != nil {
		ah.Logger.Errorw("can`t marshal wallet",
			"err:", err.Error())
		http.Error(w, "can`t make wallet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ah *WalletHandler) Update(w http.ResponseWriter, r *http.Request) {
	walletIdString := r.PathValue("walletId")
	if walletIdString == "" {
		ah.Logger.Errorw("no walletId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	walletId, err := strconv.Atoi(walletIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	wallet := &models.Wallet{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, wallet)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	//_, err = govalidator.ValidateStruct(wallet)
	//if err != nil {
	//	ah.Logger.Infow("can`t validate form",
	//		"err:", err.Error())
	//	http.Error(w, "bad data", http.StatusBadRequest)
	//	return
	//}

	wallet.ID = walletId
	err = ah.WalletUseCase.Update(wallet)
	if err != nil {
		ah.Logger.Infow("can`t update wallet",
			"err:", err.Error())
		http.Error(w, "can`t update wallet", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(wallet)

	if err != nil {
		ah.Logger.Errorw("can`t marshal wallet",
			"err:", err.Error())
		http.Error(w, "can`t make wallet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ah *WalletHandler) Delete(w http.ResponseWriter, r *http.Request) {
	walletIdString := r.PathValue("walletId")
	if walletIdString == "" {
		ah.Logger.Errorw("no walletId var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	walletId, err := strconv.Atoi(walletIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = ah.WalletUseCase.Delete(walletId)
	if err != nil {
		ah.Logger.Infow("can`t delete wallet",
			"err:", err.Error())
		http.Error(w, "can`t delete wallet", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ah *WalletHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	wallets, err := ah.WalletUseCase.GetAll()
	if err != nil {
		ah.Logger.Infow("can`t get all wallets",
			"err:", err.Error())
		http.Error(w, "can`t get all wallets", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(wallets)
	if err != nil {
		ah.Logger.Errorw("can`t marshal wallet",
			"err:", err.Error())
		http.Error(w, "can`t make wallet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

func (ah *WalletHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	walletUID := r.PathValue("WALLET_UUID")
	if walletUID == "" {
		ah.Logger.Errorw("no WALLET_UUID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	wallet, err := ah.WalletUseCase.GetByUID(walletUID)
	if err != nil {
		ah.Logger.Infow("can`t get wallet",
			"err:", err.Error())
		http.Error(w, "can`t get wallet", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(wallet)

	if err != nil {
		ah.Logger.Errorw("can`t marshal wallet",
			"err:", err.Error())
		http.Error(w, "can`t make wallet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

type ChangeAmountRequest struct {
	WalletUID     string `valid:"uuid" json:"walletId"`
	OperationType string `valid:"in(DEPOSIT|WITHDRAW)" json:"operationType"`
	Amount        int    `valid:"int" json:"amount"`
}

func (ah *WalletHandler) ChangeAmount(w http.ResponseWriter, r *http.Request) {
	changeAmountReq := ChangeAmountRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &changeAmountReq)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(changeAmountReq)
	if err != nil {
		ah.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	if changeAmountReq.Amount <= 0 {
		ah.Logger.Infow("can`t change amount",
			"err", "amount must be positive")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	amountReq := 0
	if changeAmountReq.OperationType == "DEPOSIT" {
		amountReq = changeAmountReq.Amount
	} else if changeAmountReq.OperationType == "WITHDRAW" {
		amountReq = -changeAmountReq.Amount
	}

	err = ah.WalletUseCase.ChangeAmount(changeAmountReq.WalletUID, amountReq)
	if err != nil {
		ah.Logger.Infow("can`t create wallet",
			"err:", err.Error())
		http.Error(w, "can`t create wallet", http.StatusBadRequest)
		return
	}

	//resp, err := json.Marshal(wallet)
	//
	//if err != nil {
	//	ah.Logger.Errorw("can`t marshal wallet",
	//		"err:", err.Error())
	//	http.Error(w, "can`t make wallet", http.StatusInternalServerError)
	//	return
	//}

	//w.Header().Set("Location", fmt.Sprintf("/api/v1/wallets/%d", changeAmountReq.ID))
	w.WriteHeader(http.StatusOK)

	//_, err = w.Write(resp)
	//if err != nil {
	//	ah.Logger.Errorw("can`t write response",
	//		"err:", err.Error())
	//	http.Error(w, "can`t write response", http.StatusInternalServerError)
	//	return
	//}
}
