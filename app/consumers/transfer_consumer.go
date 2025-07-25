package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"money-transfer/app/models"
	"money-transfer/app/models/constants"
	"money-transfer/app/usecases"
	"money-transfer/app/utils/helper"
	kafkautils "money-transfer/app/utils/kafka"
	"money-transfer/app/utils/sqldb"
	"net/http"
	"runtime"

	"github.com/google/uuid"
)

type TransferConsumerInterface interface {
	ProcessMessage()
}

type transferConsumer struct {
	ctx                context.Context
	sqlDBMap           map[string]sqldb.SqlInterface
	kafkaClient        kafkautils.KafkaClientInterface
	transactionUseCase usecases.TransactionUseCaseInterface
	// merchantRepository    repositories.MerchantRepositoryInterface
	// transactionRepository repositories.TransactionRepositoryInterface
}

func InitTransferConsumerHandler(
	ctx context.Context,
	sqlDBMap map[string]sqldb.SqlInterface,
	kafkaClient kafkautils.KafkaClientInterface,
	transactionUseCase usecases.TransactionUseCaseInterface,
	// merchantRepository repositories.MerchantRepositoryInterface,
	// transactionRepository repositories.TransactionRepositoryInterface,
) TransferConsumerInterface {
	return &transferConsumer{
		ctx:                ctx,
		sqlDBMap:           sqlDBMap,
		kafkaClient:        kafkaClient,
		transactionUseCase: transactionUseCase,
		// merchantRepository:    merchantRepository,
		// transactionRepository: transactionRepository,
	}
}

func (c *transferConsumer) ProcessMessage() {
	topic := constants.TRANSFER_KAFKA_TOPIC
	groupID := constants.TRANSFER_KAFKA_GROUP
	defer helper.Recover(topic)
	reader := c.kafkaClient.SetConsumerGroupReader(topic, groupID)

	numWorkers := runtime.NumCPU() * 5
	workerPool := make(chan struct{}, numWorkers)

	for {
		m, err := reader.ReadMessage(c.ctx)
		if err != nil {
			break
		}

		request := &models.TransferConsumerRequest{}
		err = json.Unmarshal(m.Value, &request)
		if err != nil {
			break
		}

		fmt.Println(request)

		workerPool <- struct{}{}
		go func(req *models.TransferConsumerRequest) {
			defer func() { <-workerPool }()
			c.ProcessTransfer(req)
		}(request)

		err = reader.CommitMessages(c.ctx, m)
		if err != nil {
			fmt.Printf("error commit kafka message [%s]\n", err.Error())
		}
	}

	if err := reader.Close(); err != nil {
		fmt.Println(err)
	}
}

func (c *transferConsumer) ProcessTransfer(request *models.TransferConsumerRequest) {
	dbTx, err := c.sqlDBMap[constants.DB_MONEY_TRANSFER].DB().BeginTx(c.ctx, nil)
	if err != nil {
		c.Retry(request)
		return
	}

	errLog := c.transactionUseCase.ProcessTransfer(dbTx, request)
	if errLog != nil {
		dbTx.Rollback()
		c.Retry(request)
		return
	}

	dbTx.Commit()
}

func (c *transferConsumer) Retry(request *models.TransferConsumerRequest) *models.ErrorLog {
	keyKafka := []byte(uuid.New().String())
	messageKafka, err := json.Marshal(request)
	if err != nil {
		err = fmt.Errorf("[PublishMessage] Error marshal request [%s]", err)
		return helper.WriteLog(err, http.StatusInternalServerError, err.Error())
	}

	err = c.kafkaClient.WriteToTopic(constants.TRANSFER_KAFKA_TOPIC, keyKafka, messageKafka)
	if err != nil {
		err = fmt.Errorf("[PublishMessage] Error WriteToTopic [%s]", err)
		return helper.WriteLog(err, http.StatusInternalServerError, err.Error())
	}

	return nil
}
