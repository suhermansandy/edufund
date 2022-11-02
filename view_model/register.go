package view_model

type Register struct {
	FullName             *string `json:"full_name"`
	UserName             *string `json:"user_name"`
	Password             *string `json:"password"`
	ConfirmationPassword *string `json:"confirmation_password"`
}
