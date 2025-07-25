package consumers

import (
	"context"
	"money-transfer/app/repositories"
	"money-transfer/app/usecases"
	kafkautils "money-transfer/app/utils/kafka"
	"money-transfer/app/utils/sqldb"
)

func InitTransferConsumer(
	ctx context.Context,
	kafkaClient kafkautils.KafkaClientInterface,
	sqlDBMap map[string]sqldb.SqlInterface,
) TransferConsumerInterface {
	merchantRepo := repositories.InitMerchantRepository(ctx, sqlDBMap)
	transactionRepo := repositories.InitTransactionRepository(ctx, sqlDBMap)
	transactionUseCase := usecases.InitTransactionUseCaseInterface(ctx, kafkaClient, merchantRepo, transactionRepo)
	return InitTransferConsumerHandler(ctx, sqlDBMap, kafkaClient, transactionUseCase)
}
