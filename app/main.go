package main

import (
	"context"
	"fmt"
	"money-transfer/app/handlers"
	"money-transfer/app/models/constants"
	kafkautils "money-transfer/app/utils/kafka"
	"money-transfer/app/utils/sqldb"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	arg := os.Args[1]

	switch arg {
	case "main":
		mainWithoutArg()
		break
	case "consumer":
		consumers()
	default:
		mainWithoutArg()
	}
}

func mainWithoutArg() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := godotenv.Load(".env")
	if err != nil {
		errStr := fmt.Sprintf(".env not load properly %s", err.Error())
		panic(errStr)
	}

	ctx := context.Background()

	moneyTransferPgsql, err := sqldb.InitPgSql("postgres", os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DATABASE"), os.Getenv("PG_SCHEMA"))
	if err != nil {
		errStr := fmt.Sprintf("Error pgsql selling out service connection %s", err.Error())
		panic(errStr)
	}

	sqlDBMap := map[string]sqldb.SqlInterface{}
	sqlDBMap[constants.DB_MONEY_TRANSFER] = moneyTransferPgsql

	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOSTS"), ",")
	kafkaClient := kafkautils.InitKafkaClientInterface(ctx, kafkaHosts)

	defer moneyTransferPgsql.DB().Close()
	defer kafkaClient.GetController().Close()
	defer kafkaClient.GetConnection().Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Printf("Starting HTTP Service\n")
		handlers.MainHttpHandler(ctx, sqlDBMap, kafkaClient)
	}()

	wg.Wait()
}

func consumers() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := godotenv.Load(".env")
	if err != nil {
		errStr := fmt.Sprintf(".env not load properly %s", err.Error())
		panic(errStr)
	}

	ctx := context.Background()

	moneyTransferPgsql, err := sqldb.InitPgSql("postgres", os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DATABASE"), os.Getenv("PG_SCHEMA"))
	if err != nil {
		errStr := fmt.Sprintf("Error pgsql selling out service connection %s", err.Error())
		panic(errStr)
	}

	sqlDBMap := map[string]sqldb.SqlInterface{}
	sqlDBMap[constants.DB_MONEY_TRANSFER] = moneyTransferPgsql

	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOSTS"), ",")
	kafkaClient := kafkautils.InitKafkaClientInterface(ctx, kafkaHosts)

	defer moneyTransferPgsql.DB().Close()
	defer kafkaClient.GetController().Close()
	defer kafkaClient.GetConnection().Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Printf("Starting Consumer Service\n")
		handlers.MainConsumerHandler(ctx, kafkaClient, sqlDBMap)
	}()

	wg.Wait()
}
