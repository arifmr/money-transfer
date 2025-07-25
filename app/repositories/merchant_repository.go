package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"money-transfer/app/models"
	"money-transfer/app/models/constants"
	"money-transfer/app/utils/helper"
	"money-transfer/app/utils/sqldb"
	"net/http"
)

type MerchantRepositoryInterface interface {
	UpdateBalanceWithLockByID(dbTx *sql.Tx, request *models.UpdateBalanceWithLockByIDRequest) *models.ErrorLog
}

type merchantRepository struct {
	ctx   context.Context
	sqlDB map[string]sqldb.SqlInterface
}

func InitMerchantRepository(ctx context.Context, sqlDB map[string]sqldb.SqlInterface) MerchantRepositoryInterface {
	return &merchantRepository{
		ctx:   ctx,
		sqlDB: sqlDB,
	}
}

func (r *merchantRepository) UpdateBalanceWithLockByID(dbTx *sql.Tx, request *models.UpdateBalanceWithLockByIDRequest) *models.ErrorLog {
	lockQuery := `
		SELECT id, balance
		FROM merchants
		WHERE id = $1
		FOR UPDATE
	`

	var merchantID int
	var currentAmount float64

	row := dbTx.QueryRowContext(r.ctx, lockQuery, request.MerchantID)
	if err := row.Scan(&merchantID, &currentAmount); err != nil {
		return helper.WriteLog(err, http.StatusInternalServerError, err.Error())
	}

	if currentAmount < request.Amount && request.Operator == constants.TRANSFER_OPERATOR {
		err := fmt.Errorf("insufficient balance: current=%.2f, required=%.2f", currentAmount, request.Amount)
		return helper.WriteLog(err, http.StatusBadRequest, err.Error())
	}

	updateQuery := fmt.Sprintf(`
		UPDATE merchants
		SET balance = balance %s $1,
		    updated_at = NOW()
		WHERE id = $2
	`, request.Operator)

	_, err := dbTx.ExecContext(r.ctx, updateQuery, request.Amount, request.MerchantID)
	if err != nil {
		return helper.WriteLog(err, http.StatusInternalServerError, err.Error())
	}

	return nil
}
