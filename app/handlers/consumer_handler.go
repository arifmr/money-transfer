package handlers

import (
	"context"
	"money-transfer/app/consumers"
	kafkautils "money-transfer/app/utils/kafka"
	"money-transfer/app/utils/sqldb"
	"sync"
)

func MainConsumerHandler(ctx context.Context, kafkaClientInterface kafkautils.KafkaClientInterface, sqlDBMap map[string]sqldb.SqlInterface) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	transferConsumer := consumers.InitTransferConsumer(ctx, kafkaClientInterface, sqlDBMap)
	go transferConsumer.ProcessMessage()

	wg.Wait()
}
