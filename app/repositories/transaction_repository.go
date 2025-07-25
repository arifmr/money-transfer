package repositories

import (
	"context"
	"database/sql"
	"money-transfer/app/models"
	"money-transfer/app/utils/helper"
	"money-transfer/app/utils/sqldb"
	"net/http"
)

type TransactionRepositoryInterface interface {
	CreateTransaction(dbTx *sql.Tx, request *models.CreateTransactionRequest) (int64, *models.ErrorLog)
	UpdateTransactionStatus(dbTx *sql.Tx, request *models.UpdateTransactionStatusRequest) *models.ErrorLog
}

type transactionRepository struct {
	ctx   context.Context
	sqlDB map[string]sqldb.SqlInterface
}

func InitTransactionRepository(ctx context.Context, sqlDB map[string]sqldb.SqlInterface) TransactionRepositoryInterface {
	return &transactionRepository{
		ctx:   ctx,
		sqlDB: sqlDB,
	}
}

func (r *transactionRepository) CreateTransaction(dbTx *sql.Tx, request *models.CreateTransactionRequest) (int64, *models.ErrorLog) {
	var id int64

	insertQuery := `
		INSERT INTO transactions (merchant_id, status, amount, account_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := dbTx.QueryRowContext(
		r.ctx,
		insertQuery,
		request.MerchantID,
		request.Status,
		request.Amount,
		request.AccountNumber,
	).Scan(&id)

	if err != nil {
		return id, helper.WriteLog(err, http.StatusInternalServerError, err.Error())
	}

	return id, nil
}

func (r *transactionRepository) UpdateTransactionStatus(dbTx *sql.Tx, request *models.UpdateTransactionStatusRequest) *models.ErrorLog {
	countQuery := `
		SELECT COUNT(*)
		FROM transactions
		WHERE id = $1
	`

	var count int
	err := dbTx.QueryRowContext(r.ctx, countQuery, request.ID).Scan(&count)
	if err != nil {
		return helper.WriteLog(err, http.StatusInternalServerError, "Failed to validate transaction ID")
	}

	if count == 0 {
		return helper.WriteLog(nil, http.StatusNotFound, "Transaction ID not found")
	}

	updateQuery := `
		UPDATE transactions
		SET status = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err = dbTx.ExecContext(r.ctx, updateQuery, request.Status, request.ID)
	if err != nil {
		return helper.WriteLog(err, http.StatusInternalServerError, err.Error())
	}

	return nil
}
