package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"money-transfer/app/models"
	"money-transfer/app/models/constants"
	"money-transfer/app/repositories"
	"money-transfer/app/utils/helper"
	kafkautils "money-transfer/app/utils/kafka"
	"net/http"

	"github.com/google/uuid"
)

type TransactionUseCaseInterface interface {
	Transfer(dbTx *sql.Tx, request *models.TransferRequest) *models.ErrorLog
	ProcessTransfer(dbTx *sql.Tx, request *models.TransferConsumerRequest) *models.ErrorLog
	MockAPICall(simulateSuccess bool) string
}

type transactionUseCase struct {
	ctx                   context.Context
	kafkaClient           kafkautils.KafkaClientInterface
	merchantRepository    repositories.MerchantRepositoryInterface
	transactionRepository repositories.TransactionRepositoryInterface
}

func InitTransactionUseCaseInterface(ctx context.Context, kafkaClient kafkautils.KafkaClientInterface, merchantRepository repositories.MerchantRepositoryInterface, transactionRepository repositories.TransactionRepositoryInterface) TransactionUseCaseInterface {
	return &transactionUseCase{
		ctx:                   ctx,
		kafkaClient:           kafkaClient,
		merchantRepository:    merchantRepository,
		transactionRepository: transactionRepository,
	}
}

func (u *transactionUseCase) Transfer(dbTx *sql.Tx, request *models.TransferRequest) *models.ErrorLog {
	updateBalanceReq := &models.UpdateBalanceWithLockByIDRequest{
		MerchantID: request.MerchantID,
		Amount:     request.Amount,
		Operator:   constants.TRANSFER_OPERATOR,
	}

	errLog := u.merchantRepository.UpdateBalanceWithLockByID(dbTx, updateBalanceReq)
	if errLog != nil {
		return errLog
	}

	createTrxReq := &models.CreateTransactionRequest{
		MerchantID:    request.MerchantID,
		Status:        constants.TRANSACTION_STATUS_PENDING,
		Amount:        request.Amount,
		AccountNumber: request.AccountNumber,
	}

	trxID, errLog := u.transactionRepository.CreateTransaction(dbTx, createTrxReq)
	if errLog != nil {
		return errLog
	}

	consumerRequest := &models.TransferConsumerRequest{
		TransactionID:   trxID,
		MerchantID:      request.MerchantID,
		Amount:          request.Amount,
		AccountNumber:   request.AccountNumber,
		SimulateSuccess: request.SimulateSuccess,
	}

	keyKafka := []byte(uuid.New().String())
	messageKafka, err := json.Marshal(consumerRequest)
	if err != nil {
		err = fmt.Errorf("[PublishMessage] Error marshal request [%s]", err)
		return helper.WriteLog(err, http.StatusInternalServerError, err.Error())
	}

	err = u.kafkaClient.WriteToTopic(constants.TRANSFER_KAFKA_TOPIC, keyKafka, messageKafka)
	if err != nil {
		err = fmt.Errorf("[PublishMessage] Error WriteToTopic [%s]", err)
		return helper.WriteLog(err, http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (u *transactionUseCase) ProcessTransfer(dbTx *sql.Tx, request *models.TransferConsumerRequest) *models.ErrorLog {
	// mock call api payment provider
	status := u.MockAPICall(request.SimulateSuccess)

	updateTrxStatusReq := &models.UpdateTransactionStatusRequest{
		ID:     request.TransactionID,
		Status: status,
	}

	errLog := u.transactionRepository.UpdateTransactionStatus(dbTx, updateTrxStatusReq)
	if errLog != nil {
		return errLog
	}

	if status == constants.TRANSACTION_STATUS_FAILED {
		refundRequest := &models.UpdateBalanceWithLockByIDRequest{
			MerchantID: request.MerchantID,
			Amount:     request.Amount,
			Operator:   constants.REFUND_OPERATOR,
		}

		errLog := u.merchantRepository.UpdateBalanceWithLockByID(dbTx, refundRequest)
		if errLog != nil {
			return errLog
		}
	}

	return nil
}

func (u *transactionUseCase) MockAPICall(simulateSuccess bool) string {
	result := constants.TRANSACTION_STATUS_SUCCESS
	if !simulateSuccess {
		result = constants.TRANSACTION_STATUS_FAILED
	}

	return result
}
