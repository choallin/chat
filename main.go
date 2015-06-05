package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = map[string]interface{}{"name": authCookie.Value}
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8088", "The addr of the application.")
	flag.Parse()
	r := newRoom()
	http.HandleFunc("/auth/callback/login", loginHandler)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	http.Handle("/orator", MustAuth(&templateHandler{filename: "orator.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	go r.run()

	log.Println("Starte den Webserver auf ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
