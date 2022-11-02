package controllers

import (
	"log"
	"net/http"

	db "edufund/repositories"
	"edufund/services"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // For gorm usage
	"github.com/rs/cors"
)

type App struct {
	Router *mux.Router
	DB     map[string]db.DBHandler
}

func (a *App) Initialize() {
	log.Println("API Init")
	a.DB = make(map[string]db.DBHandler)
	dbHandler, err := db.OpenDB("postgres", Env.DbConn)
	if err != nil {
		log.Println("Could not connect database")
		log.Fatal(err.Error())
	}
	a.DB["db"] = db.AutoMigrate(dbHandler)

	a.DB["db"].SetLogger(gorm.Logger{LogWriter: log.Default()})
	a.DB["db"].LogMode(true)

	// set router
	a.Router = mux.NewRouter()
	a.setRouters()
	log.Println("API Ready")
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	registerHandler := services.RegisterHandler{RESTHandler: services.Route(a.DB)}
	a.Post("/register", a.ihandleRequest(registerHandler, services.Registers))

	loginHandler := services.LoginHandler{RESTHandler: services.Route(a.DB)}
	a.Post("/login", a.ihandleRequest(loginHandler, services.Login))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"*"},
		ExposedHeaders:     []string{"File-Name"},
		OptionsPassthrough: false,
	})
	log.Fatal(http.ListenAndServe(host, c.Handler(a.Router)))
}

// RequestIHandlerFunction is abbreviation for handler data, response write, request
type RequestIHandlerFunction func(h services.RESTHandler, w http.ResponseWriter, r *http.Request)

func (a *App) ihandleRequest(h services.RESTHandler, handler RequestIHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(h, w, r)
	}
}
