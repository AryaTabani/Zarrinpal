package models

type Payment struct {
	ID          int64
	Amount      int
	Description string
	Authority   string
	RefID       string
	Status      string
}

type PaymentRequestPayload struct {
	Amount      int      `json:"amount" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Metadata    Metadata `json:"metadata"`
}

type Metadata struct {
	Mobile string `json:"mobile,omitempty"`
	Email  string `json:"email,omitempty"`
}

type ZarinpalRequest struct {
	MerchantID  string   `json:"merchant_id"`
	Amount      int      `json:"amount"`
	Description string   `json:"description"`
	CallbackURL string   `json:"callback_url"`
	Metadata    Metadata `json:"metadata,omitempty"`
}

type ZarinpalResponse struct {
	Data struct {
		Authority string `json:"authority"`
		Code      int    `json:"code"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

type ZarinpalVerifyRequest struct {
	MerchantID string `json:"merchant_id"`
	Amount     int    `json:"amount"`
	Authority  string `json:"authority"`
}

type ZarinpalVerifyResponse struct {
	Data struct {
		Code  int   `json:"code"`
		RefID int64 `json:"ref_id"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}
