package handlers

import (
	"context"
	"fmt"
	"money-transfer/app/routes"
	kafkautils "money-transfer/app/utils/kafka"
	"money-transfer/app/utils/sqldb"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MainHttpHandler(ctx context.Context, sqlDBMap map[string]sqldb.SqlInterface, kafkaClient kafkautils.KafkaClientInterface) {
	g := gin.Default()
	g.Use(RequestId())
	routes.InitHTTPRoute(ctx, g, sqlDBMap, kafkaClient)
	addr := fmt.Sprintf(":%s", os.Getenv("MAIN_PORT"))
	http.ListenAndServe(addr, g)
}

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get("X-Request-Id")

		// Create request id with UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Expose it for use in the application
		c.Set("RequestId", requestID)
		// Set X-Request-Id header
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}
