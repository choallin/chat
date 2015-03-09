package main

import "net/http"

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
	if r.PostFormValue("login") == "stage" && r.PostFormValue("password") == "nico" {
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: "stage",
			Path:  "/",
		})
		//@todo: JSON output, so that JavaScript can handle the state of the UI

	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
