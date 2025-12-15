package gettransactions

type TransactionModel struct {
	TransactionId string  `json:"transaction_id"`
	AccountId     string  `json:"account_id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Note          string  `json:"note"`
	CreatedAt     string  `json:"created_at"`
}
