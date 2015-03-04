package main

import (
	"log"
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (auth *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		//User logged in. Call next handler
		auth.next.ServeHTTP(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("%v", r.PostFormValue("login"))
	log.Println("%v", r.PostFormValue("password"))
	//check if login Data is Valid
	//connect to Redis Server that has the user login data:
	//
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
