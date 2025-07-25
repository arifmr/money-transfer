package controllers

import (
	"context"
	"money-transfer/app/repositories"
	"money-transfer/app/usecases"
	kafkautils "money-transfer/app/utils/kafka"
	"money-transfer/app/utils/sqldb"
)

func InitHTTPTransactionController(ctx context.Context, sqlDBMap map[string]sqldb.SqlInterface, kafkaClient kafkautils.KafkaClientInterface) TransactionControllerInterface {
	transactionRepository := repositories.InitTransactionRepository(ctx, sqlDBMap)
	merchantRepository := repositories.InitMerchantRepository(ctx, sqlDBMap)
	transactionUseCase := usecases.InitTransactionUseCaseInterface(ctx, kafkaClient, merchantRepository, transactionRepository)
	transactionValidator := usecases.InitTransactionValidator()
	return InitTransactionController(ctx, sqlDBMap, transactionUseCase, transactionValidator)
}
