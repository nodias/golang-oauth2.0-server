package router

import (
	"github.com/go-session/session"
	"github.com/nodias/golang-oauth2.0-common/models"
	"log"
	"net/http"
	"net/url"
)

func authHandler(w http.ResponseWriter, req *http.Request) {

	store, err := session.Start(nil, w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		//w.Header().Set("Location", "/login")
		err := renderer.HTML(w, http.StatusFound, "login", &models.Response{
			Id:    "",
			User:  nil,
			Error: nil,
		})
		if err != nil {
			log.Println(err)
			return
		}

		return
	}

	err = renderer.HTML(w, http.StatusOK, "auth", &models.Response{
		Id:    "",
		User:  nil,
		Error: nil,
	})
	if err != nil {
		log.Println(err)
		return
	}

}

func tokenHandler(w http.ResponseWriter, req *http.Request) {
	err := srv.HandleTokenRequest(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	r.Form = form

	store.Delete("ReturnUri")
	store.Save()

	// Authorization 요청에 대한 Validation Check
	// User ID 추출 만약 인증 되지 않다면, login 화면으로 Redirect
	// Authorization Code 생성
	err = srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}
