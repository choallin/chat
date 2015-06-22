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
	if (r.PostFormValue("login") == "stage" || r.PostFormValue("login") == "regie") && r.PostFormValue("password") == "nico" {
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: r.PostFormValue("login"),
			Path:  "/",
		})
		resp := []byte("{\"valid\": true, \"location\": \"chat\"}")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	} else if r.PostFormValue("login") == "pult" && r.PostFormValue("password") == "nico" {
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: r.PostFormValue("login"),
			Path:  "/",
		})
		resp := []byte("{\"valid\": true, \"location\": \"pult\"}")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	}
	resp := []byte("{\"valid\": false}")
	w.Write(resp)
}

// MustAuth Authentification Handler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
