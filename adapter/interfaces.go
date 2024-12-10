package adapter

import "scarlet_backend/model"

type UserService interface {
	SaveByEmail(user *model.User) (*model.User, error)
	FindAll() ([]model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindById(id int) (*model.User, error)
	CheckLogin(email string, psw string) (*model.User, error)
	SendOTP(phone string) (string, error)
	VerifyOTP(otp_id string, otp_code string) error
	FindByPhone(phone string) (*model.User, error)
	SaveByPhone(user *model.User) (*model.User, error)
}
