package models

type VerifyEmail struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
