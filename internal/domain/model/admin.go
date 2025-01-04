package model

type CreateStaffRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type CreateStaffResponse struct {
	Id int64 `json:"id"`
}

type UpdateStaffRequest struct {
	Id          int64  `json:"id" binding:"required"` // required
	Name        string `json:"name"`
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type GetExternalTransactionRequest struct {
	BankId   int64  `json:"bankId"`
	FromDate string `json:"fromDate" binding:"required"` // yyyy-mm-dd
	ToDate   string `json:"toDate" binding:"required"`   // yyyy-mm-dd
}
