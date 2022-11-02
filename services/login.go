package services

import (
	"edufund/repositories"
	"edufund/util"
	"edufund/view_model"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type LoginHandler struct {
	RESTHandler
}

func (h LoginHandler) New() interface{} {
	return &view_model.Login{}
}

func (h LoginHandler) Validate(state string, v interface{}) error {
	result, ok := v.(*view_model.Login)
	if !ok {
		return errors.New("invalid type")
	}

	if _, err := util.ValidMailAddress(*result.UserName); err != nil {
		return errors.New("please provide a valid email address")
	}

	if len(*result.Password) < 12 {
		return errors.New("password should be at least 12 characters long")
	}

	status := repositories.CheckUserPass(h.DefaultDB(state), result.UserName, result.Password)
	if status == 0 {
		return errors.New("invalid username / password")
	}

	return nil
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (h LoginHandler) CreatedJwt(v interface{}) string {
	result, _ := v.(*view_model.Login)
	expirationTime := time.Now().Add(5 * time.Minute)
	var jwtKey = []byte("sandy")
	claim := &Claims{
		Username: *result.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}
