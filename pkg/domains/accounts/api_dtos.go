package accounts

type AccountAPIResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

func NewAPIResponseFromEntity(account *Account) *AccountAPIResponse {
	return &AccountAPIResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	}
}

type CreateAccountAPIRequest struct {
	DocumentNumber string `json:"document_number" validate:"required,numeric"`
}
