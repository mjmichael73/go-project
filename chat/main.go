package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(
			template.ParseFiles(
				filepath.Join("templates", t.filename),
			),
		)
		t.templ.Execute(w, r)
	})
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	gomniauth.SetSecurityKey("abcdefghijklmnopqrstuvwxyz")
	gomniauth.WithProviders(
		google.New("708077468376-27erjfilu5g33p2mdnkncg77jn33uae3.apps.googleusercontent.com", "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", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
