package services

import (
	db "edufund/repositories"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
)

type RESTHandler interface {
	DefaultDB(state string) db.DBHandler
	DB(state string) db.DBHandler
	New() interface{}
	NewRegister() interface{}
	Validate(state string, v interface{}) error
	CompareModelRegister(x interface{}, y interface{}) interface{}
	CreatedJwt(x interface{}) string
}

type DefaultHandler struct {
	DBHandler map[string]db.DBHandler
}

func (h DefaultHandler) DefaultDB(state string) db.DBHandler {
	return h.DBHandler[state]
}

func (h DefaultHandler) DB(state string) db.DBHandler {
	return h.DBHandler[state]
}

func (h DefaultHandler) New() interface{} {
	return gorm.Model{}
}

func (h DefaultHandler) NewRegister() interface{} {
	return gorm.Model{}
}

func (h DefaultHandler) Validate(state string, v interface{}) error {
	return nil
}

func (h DefaultHandler) CompareModelRegister(x interface{}, y interface{}) interface{} {
	return make([]gorm.Model, 0)
}

func (h DefaultHandler) CreatedJwt(x interface{}) string {
	return ""
}

func Route(db map[string]db.DBHandler) RESTHandler {
	return DefaultHandler{DBHandler: db}
}

func Registers(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	state := "db"

	result := h.NewRegister()
	resultLogin := h.New()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.Validate(state, result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resultFix := h.CompareModelRegister(result, resultLogin)

	tx := h.DB(state).Begin()
	if err := tx.Save(resultFix).Error(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()

	respondJSON(w, http.StatusCreated, resultFix)
}

func Login(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	var path = "history"
	if state == "mock" {
		path += "/mock"
	} else {
		state = "db"
	}

	result := h.New()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.Validate(state, result); err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token := h.CreatedJwt(result)

	respondJSON(w, http.StatusOK, token)
}
