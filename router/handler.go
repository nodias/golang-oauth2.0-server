package router

import (
	"encoding/json"
	"github.com/go-session/session"
	"github.com/gorilla/mux"
	"github.com/nodias/golang-oauth2.0-common/models"
	"github.com/nodias/golang-oauth2.0-common/shared/logger"
	"github.com/nodias/golang-oauth2.0-server/service"
	"go.elastic.co/apm"
	"log"
	"net/http"
	"strings"
	"unicode"
)

//templates
func loginHandler(w http.ResponseWriter, req *http.Request) {
	err := renderer.HTML(w, http.StatusOK, "login", &models.Response{
		Id:    "",
		User:  nil,
		Error: nil,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	store.Set("LoggedInUserID", "000000")
	err = store.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Location", "/auth")
	w.WriteHeader(http.StatusFound)

	err = renderer.HTML(w, http.StatusOK, "login", &models.Response{
		Id:    "",
		User:  nil,
		Error: nil,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

//getUserInfoHandler is a function, gets the information of one User
func getUserInfoHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.New(ctx)

	id := mux.Vars(req)["id"]
	log.WithField("id", id).Debug("handling hello request")
	if strings.IndexFunc(id, func(r rune) bool { return r >= unicode.MaxASCII }) >= 0 {
		panic("non-ASCII id!")
	}

	user, rerr := service.GetUserInfo(req.Context(), id)
	if rerr != nil {
		w.WriteHeader(rerr.Code)
	}
	err := json.NewEncoder(w).Encode(models.Response{
		Id:    models.ID(id),
		User:  user,
		Error: rerr,
	})
	if err != nil {
		rerr2 := models.NewResponseError(err, 500)
		log.WithError(rerr2).Error("failed to GetUserInfoHandler")
		//apm server에 에러를 업로드 시켜줍니다.
		apm.CaptureError(ctx, rerr2.Err).Send()
		http.Error(w, "failed encode to json", 9999)
		return
	}
}
