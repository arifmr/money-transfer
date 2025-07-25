package routes

import (
	"context"
	"money-transfer/app/controllers"
	kafkautils "money-transfer/app/utils/kafka"
	"money-transfer/app/utils/sqldb"

	"github.com/gin-gonic/gin"
)

func InitHTTPRoute(ctx context.Context, g *gin.Engine, sqlDBMap map[string]sqldb.SqlInterface, kafkaClient kafkautils.KafkaClientInterface) {
	transactionController := controllers.InitHTTPTransactionController(ctx, sqlDBMap, kafkaClient)

	g.POST("/transfer", transactionController.Transfer)
}
