package usecases

import (
	"errors"
	"money-transfer/app/models"
	"money-transfer/app/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionValidatorInterface interface {
	TransferValidator(ginCtx *gin.Context) (*models.TransferRequest, *models.ErrorLog)
}

type transactionValidator struct {
}

func InitTransactionValidator() TransactionValidatorInterface {
	return &transactionValidator{}
}

func (u *transactionValidator) TransferValidator(ginCtx *gin.Context) (*models.TransferRequest, *models.ErrorLog) {
	request := &models.TransferRequest{}

	err := ginCtx.Bind(request)
	if err != nil {
		err = errors.New("Format JSON tidak sesuai, silahkan periksa kembali")
		return nil, helper.WriteLog(err, http.StatusBadRequest, err.Error())
	}

	if request.Amount <= 0 {
		err := errors.New("amount tidak boleh kurang dari sama dengan 0")
		return request, helper.WriteLog(err, http.StatusBadRequest, nil)
	}

	if request.MerchantID == 0 {
		err := errors.New("merchant_id harus diisi dan tidak boleh 0")
		return request, helper.WriteLog(err, http.StatusBadRequest, nil)
	}

	if request.AccountNumber <= 0 {
		err := errors.New("account_number tidak valid")
		return request, helper.WriteLog(err, http.StatusBadRequest, nil)
	}

	return request, nil
}
