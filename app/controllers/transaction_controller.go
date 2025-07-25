package controllers

import (
	"context"
	"money-transfer/app/models"
	"money-transfer/app/models/constants"
	"money-transfer/app/usecases"
	"money-transfer/app/utils/helper"
	"money-transfer/app/utils/sqldb"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionControllerInterface interface {
	Transfer(ctx *gin.Context)
}

type transactionController struct {
	ctx                  context.Context
	sqlDBMap             map[string]sqldb.SqlInterface
	transactionUseCase   usecases.TransactionUseCaseInterface
	transactionValidator usecases.TransactionValidatorInterface
}

func InitTransactionController(
	ctx context.Context,
	sqlDBMap map[string]sqldb.SqlInterface,
	transactionUseCase usecases.TransactionUseCaseInterface,
	transactionValidator usecases.TransactionValidatorInterface,
) TransactionControllerInterface {
	return &transactionController{
		ctx:                  ctx,
		sqlDBMap:             sqlDBMap,
		transactionUseCase:   transactionUseCase,
		transactionValidator: transactionValidator,
	}
}

func (c *transactionController) Transfer(ginCtx *gin.Context) {
	result := models.Response{}

	request, errLog := c.transactionValidator.TransferValidator(ginCtx)
	if errLog != nil {
		result.StatusCode = errLog.StatusCode
		result.Error = errLog
		ginCtx.JSON(result.StatusCode, result)
		return
	}

	dbTx, err := c.sqlDBMap[constants.DB_MONEY_TRANSFER].DB().BeginTx(c.ctx, nil)
	if err != nil {
		result.StatusCode = http.StatusInternalServerError
		result.Error = helper.WriteLog(err, http.StatusInternalServerError, nil)
		ginCtx.JSON(result.StatusCode, result)
		return
	}

	errLog = c.transactionUseCase.Transfer(dbTx, request)
	if errLog != nil {
		dbTx.Rollback()
		result.StatusCode = errLog.StatusCode
		result.Error = errLog
		ginCtx.JSON(result.StatusCode, result)
		return
	}

	dbTx.Commit()

	result.StatusCode = http.StatusOK
	ginCtx.JSON(http.StatusOK, result)
}
