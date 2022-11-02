package repositories

import (
	"edufund/model"
	"edufund/util"
)

func CheckUserPass(dbHandler DBHandler, username *string, password *string) int {
	if username == nil || password == nil {
		return 0
	}
	encryptPass := util.Encrypt(*password)
	db := dbHandler.Model(model.User{})
	model := model.User{}
	db.Where("user_name = ? AND password = ?", username, &encryptPass).Last(&model)
	if model.UserName == nil {
		return 0
	}
	return 1
}
