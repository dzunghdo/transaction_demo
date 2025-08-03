package dto

type TransactionDTO struct {
	SourceAccountID      uint64  `json:"source_account_id" validate:"required,number,gt=0"`
	DestinationAccountID uint64  `json:"destination_account_id" validate:"required,number,gt=0"`
	Amount               float64 `json:"amount" validate:"required,number,gt=0"`
}

// Validate validates the TransactionDTO struct.
func (t TransactionDTO) Validate() error {
	return GetValidator().Struct(t)
}
