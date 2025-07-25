package models

type TransferRequest struct {
	MerchantID      int64   `json:"merchant_id"`
	Amount          float64 `json:"amount"`
	AccountNumber   int64   `json:"account_number"`
	SimulateSuccess bool    `json:"simulate_success"`
}

type CreateTransactionRequest struct {
	MerchantID    int64   `json:"merchant_id"`
	Amount        float64 `json:"amount"`
	AccountNumber int64   `json:"account_number"`
	Status        string  `json:"status"`
}

type UpdateTransactionStatusRequest struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type TransferConsumerRequest struct {
	TransactionID   int64   `json:"transaction_id"`
	MerchantID      int64   `json:"merchant_id"`
	Amount          float64 `json:"amount"`
	AccountNumber   int64   `json:"account_number"`
	SimulateSuccess bool    `json:"simulate_success"`
}

type UpdateBalanceWithLockByIDRequest struct {
	MerchantID int64   `json:"merchant_id"`
	Amount     float64 `json:"amount"`
	Operator   string  `json:"operator"`
}
