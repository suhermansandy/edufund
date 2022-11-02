package services

import (
	"edufund/model"
	"edufund/util"
	"edufund/view_model"
	"errors"
)

type RegisterHandler struct {
	RESTHandler
}

func (h RegisterHandler) NewRegister() interface{} {
	return &view_model.Register{}
}

func (h RegisterHandler) New() interface{} {
	return &model.User{}
}

func (h RegisterHandler) Validate(state string, v interface{}) error {
	result, ok := v.(*view_model.Register)
	if !ok {
		return errors.New("invalid type")
	}

	if len(*result.FullName) < 2 {
		return errors.New("name should be 2 characters or more")
	}

	if len(*result.Password) < 12 {
		return errors.New("password should be at least 12 characters long")
	}

	if _, err := util.ValidMailAddress(*result.UserName); err != nil {
		return errors.New("please provide a valid email address")
	}

	if *result.Password != *result.ConfirmationPassword {
		return errors.New("confirmation password does not match")
	}

	encryptPass := util.Encrypt(*result.Password)
	result.Password = &encryptPass

	return nil
}

func (h RegisterHandler) CompareModelRegister(x interface{}, y interface{}) interface{} {
	resultRegister, ok := x.(*view_model.Register)
	if !ok {
		return errors.New("invalid type")
	}
	resultUser, ok := y.(*model.User)
	if !ok {
		return errors.New("invalid type")
	}
	resultUser.FullName = resultRegister.FullName
	resultUser.UserName = resultRegister.UserName
	resultUser.Password = resultRegister.Password

	return &resultUser
}
