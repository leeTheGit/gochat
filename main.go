package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (this *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.once.Do(func() {
		this.templ = template.Must(
			template.ParseFiles(
				filepath.Join("templates", this.filename),
			),
		)
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	authCookie, err := r.Cookie("auth")
	if err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	this.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags

	gomniauth.SetSecurityKey("lsahjlkfjsadklfhjklasjfkhdfkljahsfdsaf7934675jhgkaf7y34q2glsa42352trewgrew")
	gomniauth.WithProviders(
		google.New("343862809262-3p0258f5aioph518qd7k7mj0rjek4q9n.apps.googleusercontent.com", "fIzYet8xMykIObTUFFLCiSxN", "http://localhost:8081/auth/callback/google"),
	)

	r := newRoom()

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout/", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// 343862809262-3p0258f5aioph518qd7k7mj0rjek4q9n.apps.googleusercontent.com
// fIzYet8xMykIObTUFFLCiSxN
