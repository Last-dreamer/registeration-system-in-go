package models

type UserRegister struct {
	UserName string `json:"username"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
