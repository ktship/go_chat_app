package main

import (
	"log"
	"net/http"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"github.com/ktship/trace"
	"os"
//	"github.com/stretchr/gomniauth/providers/facebook"
//	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

func main() {
	var addr = flag.String("addr", ":8080", "The addr of ther application.")
	flag.Parse()

	// Set up gomniauth
	gomniauth.SetSecurityKey("ktship")
	gomniauth.WithProviders(
		google.New("736876228536-g6shkrdsjl3b7uhapuh1i69s7522jmi1.apps.googleusercontent.com", "E__Zw3XUf_PY_CYKceU8r6b9", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{ filename: "chat.html" }))
	http.Handle("/login", &templateHandler{ filename: "login.html" })
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()

	log.Println("Starting web server on", *addr)

	if err := http.ListenAndServe(*addr, nil) ; err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type templateHandler struct {
	once 		sync.Once
	filename 	string
	templ 		*template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates\\", t.filename)))
	})
	data := map[string]interface{} {
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}
