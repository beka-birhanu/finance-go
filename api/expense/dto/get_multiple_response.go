package dto

type GetMultipleResponse struct {
	Expenses []*GetExpenseResponse `json:"expenses"`
	Cursor   string                `json:"cursor"`
}
