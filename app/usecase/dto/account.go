package dto

type AccountDTO struct {
	AccountID uint64  `json:"account_id" validate:"required,number,gt=0"`
	Balance   float64 `json:"balance" validate:"required,number,gt=0"`
}

// Validate validates the AccountDTO struct.
func (a AccountDTO) Validate() error {
	return GetValidator().Struct(a)
}
