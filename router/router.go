package router

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"go.elastic.co/apm/module/apmgorilla"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	oauthModels "gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"log"
	"net/http"
)

const (
	//fileServer
	fsPathPrefix = "/static/"
	fsDir        = "web/"

	//html
	htmlPathPrefix = "/"

	//userApi
	userApiPathPrefix = "/api/v1/users/"
	idPattern         = "/{id:[0-9]+}"
)

var renderer *render.Render
var templatesPath = fmt.Sprintf("%s/templates", fsDir)

var manager = manage.NewDefaultManager()
var srv = server.NewServer(server.NewConfig(), manager)

func init() {
	renderer = render.New(render.Options{Directory: templatesPath})
}

func NewRouter() *mux.Router {
	return router()
}

func router() *mux.Router {
	r := mux.NewRouter()

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))
	clientStore := store.NewClientStore()
	clientStore.Set("kss", &oauthModels.Client{
		ID:     "kss",
		Secret: "secretkss",
		Domain: "http://localhost:7012",
	})
	manager.MapClientStorage(clientStore)
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	//ROUTER
	//file server router
	r.PathPrefix(fsPathPrefix).Handler(http.StripPrefix("/", http.FileServer(http.Dir(fsDir))))

	//template router
	tr := r.PathPrefix(htmlPathPrefix).Subrouter()
	tr.HandleFunc("/login", loginHandler).Methods("GET")

	//oauth2
	tr.HandleFunc("/login", loginPostHandler).Methods("POST")
	tr.HandleFunc("/auth", authHandler).Methods("GET")
	tr.HandleFunc("/token", tokenHandler).Methods("POST")
	tr.HandleFunc("/authorize", authorizeHandler)

	//api router
	ar := r.PathPrefix(userApiPathPrefix).Subrouter()
	ar.HandleFunc(idPattern, getUserInfoHandler).Methods("GET")
	//ar.HandleFunc(idPattern, getUserInfoHandler).Methods("PUT")
	//ar.HandleFunc("/", getUserInfoHandler).Methods("POST")
	//ar.HandleFunc(idPattern, getUserInfoHandler).Methods("DELETE")
	ar.Use(apmgorilla.Middleware())
	return r
}
